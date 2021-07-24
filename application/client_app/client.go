package client_app

import (
	"context"
	"time"

	"go.mercari.io/datastore/boom"

	"gae-go-sample/adapter"
	pb "gae-go-sample/proto/go/pb"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	"gae-go-sample/domain"
)

type clientApplication struct {
	executorID domain.ClientID
	*application
}

func (a *clientApplication) GetAll(ctx context.Context, page int32, offset int32) ([]*domain.Client, error) {
	me, err := a.clientRepository.Get(ctx, a.executorID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	clients, err := a.clientRepository.GetAllByCompanyWithPager(ctx, me.CompanyID, domain.NewPager(page, offset))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for i := range clients {
		urlWithSignature, err := a.publishResourceService(ctx, clients[i].GSIconURL)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		clients[i].SignedIconURL = urlWithSignature
	}

	return a.bindClientEmail(ctx, clients)
}

func (a *clientApplication) GetAllByIDs(ctx context.Context, ids []domain.ClientID) ([]*domain.Client, error) {
	clients, err := a.clientRepository.GetMulti(ctx, ids)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for i := range clients {
		urlWithSignature, err := a.publishResourceService(ctx, clients[i].GSIconURL)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		clients[i].SignedIconURL = urlWithSignature
	}

	return a.bindClientEmail(ctx, clients)
}

func (a *clientApplication) Create(
	ctx context.Context,
	email string,
	password string,
	params adapter.CreateClientParams,
	now time.Time) (*domain.Client, error) {
	me, err := a.clientRepository.Get(ctx, a.executorID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if !me.IsAdmin() {
		return nil, errors.WithStack(domain.ErrForbiddenClientRole)
	}

	company, err := a.companyRepository.Get(ctx, me.CompanyID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	fireUser, err := a.fireUserRepository.Create(ctx, email, password)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	user := domain.FromFireUser(fireUser, pb.User_Role_Client)

	client, err := domain.NewClient(
		user,
		company.ID,
		params.Name,
		params.NameKana,
		params.PhoneNumber,
		params.Position,
		params.Role,
		now)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err := a.transaction(ctx, func(tx *boom.Transaction) error {
		if err := a.userRepository.Put(tx, user); err != nil {
			return err
		}

		if err := a.clientRepository.Put(tx, client); err != nil {
			return err
		}

		return nil
	}); err != nil {
		a.rollbackFireUserService(ctx, fireUser.UID)
		return nil, errors.WithStack(err)
	}

	rooms, err := a.messageRoomRepository.GetAllByCompany(ctx, me.CompanyID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	eg := errgroup.Group{}

	for i := range rooms {
		room := rooms[i]

		eg.Go(func() error {
			if err := a.rtMemberRepository.UpdateClient(ctx, room, client); err != nil {
				return err
			}

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, errors.WithStack(err)
	}

	urlWithSignature, err := a.publishResourceService(ctx, me.GSIconURL)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	me.SignedIconURL = urlWithSignature

	return client, nil
}

func (a *clientApplication) Update(
	ctx context.Context,
	params adapter.UpdateClientParams) (*domain.Client, error) {
	me, err := a.clientRepository.Get(ctx, a.executorID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err := me.Update(
		params.Name,
		params.NameKana,
		params.IconURL,
		params.PhoneNumber,
		params.Position); err != nil {
		return nil, errors.WithStack(err)
	}

	if err := a.transaction(ctx, func(tx *boom.Transaction) error {
		if err := a.clientRepository.Put(tx, me); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, errors.WithStack(err)
	}

	rooms, err := a.messageRoomRepository.GetAllByCompany(ctx, me.CompanyID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	eg := errgroup.Group{}

	for i := range rooms {
		room := rooms[i]

		eg.Go(func() error {
			if err := a.rtMemberRepository.UpdateClient(ctx, room, me); err != nil {
				return err
			}

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, errors.WithStack(err)
	}

	urlWithSignature, err := a.publishResourceService(ctx, me.GSIconURL)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	me.SignedIconURL = urlWithSignature

	return me, nil
}

func (a *clientApplication) UpdateRole(
	ctx context.Context,
	id domain.ClientID,
	role pb.Client_Role) (*domain.Client, error) {
	me, err := a.clientRepository.Get(ctx, a.executorID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if !me.IsAdmin() {
		return nil, errors.WithStack(domain.ErrForbiddenClientRole)
	}

	client, err := a.clientRepository.Get(ctx, id)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if client.CompanyID != me.CompanyID {
		return nil, errors.WithStack(domain.ErrInvalidClient)
	}

	adminCount, err := a.clientRepository.GetCountByCompanyAndRole(ctx, me.CompanyID, pb.Client_Role_Admin)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if role == pb.Client_Role_Member && adminCount <= 1 {
		return nil, errors.WithStack(domain.ErrForbiddenClientRole)
	}

	client.UpdateRole(role)

	if err := a.transaction(ctx, func(tx *boom.Transaction) error {
		if err := a.clientRepository.Put(tx, client); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, errors.WithStack(err)
	}

	urlWithSignature, err := a.publishResourceService(ctx, client.GSIconURL)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	client.SignedIconURL = urlWithSignature

	return client, nil
}

func (a *clientApplication) Delete(
	ctx context.Context,
	id domain.ClientID) error {
	me, err := a.clientRepository.Get(ctx, a.executorID)
	if err != nil {
		return errors.WithStack(err)
	}

	if !me.IsAdmin() {
		return errors.WithStack(domain.ErrForbiddenClientRole)
	}

	client, err := a.clientRepository.Get(ctx, id)
	if err != nil {
		return errors.WithStack(err)
	}

	if client.CompanyID != me.CompanyID {
		return errors.WithStack(domain.ErrInvalidClient)
	}

	adminCount, err := a.clientRepository.GetCountByCompanyAndRole(ctx, me.CompanyID, pb.Client_Role_Admin)
	if err != nil {
		return errors.WithStack(err)
	}

	if client.IsAdmin() && adminCount <= 1 {
		return errors.WithStack(domain.ErrForbiddenClientRole)
	}

	client.Delete()

	if err := a.transaction(ctx, func(tx *boom.Transaction) error {
		if err := a.clientRepository.Put(tx, client); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

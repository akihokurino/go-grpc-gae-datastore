package customer_app

import (
	"context"
	"time"

	"go.mercari.io/datastore/boom"

	pb "gae-go-recruiting-server/proto/go/pb"

	"github.com/pkg/errors"

	"gae-go-recruiting-server/adapter"
	"gae-go-recruiting-server/domain"

	"golang.org/x/sync/errgroup"
)

type customerApplication struct {
	executorID domain.CustomerID
	*application
}

func (a *customerApplication) Create(
	ctx context.Context,
	name string,
	nameKana string,
	gender pb.User_Gender,
	phoneNumber string,
	birthdate time.Time,
	now time.Time) (*domain.Customer, error) {
	isExists, err := a.userRepository.Exists(ctx, domain.UserID(a.executorID))
	if err != nil {
		a.rollbackFireUserService(ctx, domain.UserID(a.executorID))
		return nil, errors.WithStack(err)
	}

	if isExists {
		return nil, errors.WithStack(domain.ErrUserAlreadyExists)
	}

	fireUser, err := a.fireUserRepository.Get(ctx, domain.UserID(a.executorID))
	if err != nil {
		a.rollbackFireUserService(ctx, domain.UserID(a.executorID))
		return nil, errors.WithStack(err)
	}

	user := domain.FromFireUser(fireUser, pb.User_Role_Customer)

	customer, err := domain.NewDefaultCustomer(
		user,
		name,
		nameKana,
		gender,
		phoneNumber,
		birthdate,
		now)
	if err != nil {
		a.rollbackFireUserService(ctx, domain.UserID(a.executorID))
		return nil, errors.WithStack(err)
	}

	if err := a.transaction(ctx, func(tx *boom.Transaction) error {
		if err := a.userRepository.Put(tx, user); err != nil {
			return err
		}

		if err := a.customerRepository.Put(tx, customer); err != nil {
			return err
		}

		return nil
	}); err != nil {
		a.rollbackFireUserService(ctx, domain.UserID(a.executorID))
		return nil, errors.WithStack(err)
	}

	if err := a.customerIndexRepository.Put(ctx, customer); err != nil {
		return nil, errors.WithStack(err)
	}

	urlWithSignature, err := a.publishResourceService(ctx, customer.GSIconURL)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	customer.SignedIconURL = urlWithSignature

	return customer, nil
}

func (a *customerApplication) Update(
	ctx context.Context,
	params adapter.CustomerParams) (*domain.Customer, error) {
	me, err := a.customerRepository.Get(ctx, a.executorID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err := me.Update(
		params.Name,
		params.NameKana,
		params.IconURL,
		params.Birthdate,
		params.Gender,
		params.PhoneNumber,
		params.Pr,
		params.Address,
		params.ResumeURL); err != nil {
		return nil, errors.WithStack(err)
	}

	if err := a.transaction(ctx, func(tx *boom.Transaction) error {
		if err := a.customerRepository.Put(tx, me); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, errors.WithStack(err)
	}

	if err := a.customerIndexRepository.Put(ctx, me); err != nil {
		return nil, errors.WithStack(err)
	}

	rooms, err := a.messageRoomRepository.GetAllByCustomer(ctx, me.ID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	eg := errgroup.Group{}

	for i := range rooms {
		room := rooms[i]

		eg.Go(func() error {
			if err := a.rtMemberRepository.UpdateCustomer(ctx, room, me); err != nil {
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

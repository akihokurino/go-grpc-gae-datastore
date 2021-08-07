package apply_client_app

import (
	"context"
	"time"

	"go.mercari.io/datastore/boom"

	"gae-go-recruiting-server/domain"
	pb "gae-go-recruiting-server/proto/go/pb"

	"github.com/pkg/errors"
)

type adminApplication struct {
	executorID domain.AdminUserID
	*application
}

func (a *adminApplication) GetAllWithPager(ctx context.Context, page int32, offset int32) ([]*domain.ApplyClient, error) {
	applies, err := a.applyClientRepository.GetAllWithPager(ctx, domain.NewPager(page, offset))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return applies, nil
}

func (a *adminApplication) GetAllByFilterWithPager(
	ctx context.Context,
	status pb.ApplyClient_Status,
	page int32,
	offset int32) ([]*domain.ApplyClient, error) {
	applies, err := a.applyClientRepository.GetAllByStatusWithPager(ctx, status, domain.NewPager(page, offset))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return applies, nil
}

func (a *adminApplication) Get(ctx context.Context, email domain.ApplyClientID) (*domain.ApplyClient, error) {
	apply, err := a.applyClientRepository.Get(ctx, email)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return apply, nil
}

func (a *adminApplication) GetTotalCount(ctx context.Context) (int64, error) {
	count, err := a.applyClientRepository.GetTotalCount(ctx)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return count, nil
}

func (a *adminApplication) GetCountByFilter(ctx context.Context, status pb.ApplyClient_Status) (int64, error) {
	count, err := a.applyClientRepository.GetCountByStatus(ctx, status)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return count, nil
}

func (a *adminApplication) Accept(
	ctx context.Context,
	email domain.ApplyClientID,
	password string,
	now time.Time) error {
	apply, err := a.applyClientRepository.Get(ctx, email)
	if err != nil {
		return errors.WithStack(err)
	}

	apply.Accept()

	fireUser, err := a.fireUserRepository.Create(ctx, apply.Email.String(), password)
	if err != nil {
		return errors.WithStack(err)
	}

	user := domain.FromFireUser(fireUser, pb.User_Role_Client)

	company, client, err := apply.CreateCompanyWithClient(
		domain.CompanyID(a.idFactory.UUID()),
		user,
		now,
	)
	if err != nil {
		a.rollbackFireUserService(ctx, fireUser.UID)
		return errors.WithStack(err)
	}

	if err := a.transaction(ctx, func(tx *boom.Transaction) error {
		if err := a.userRepository.Put(tx, user); err != nil {
			return err
		}

		if err := a.companyRepository.Put(tx, company); err != nil {
			return err
		}

		if err := a.clientRepository.Put(tx, client); err != nil {
			return err
		}

		if err := a.applyClientRepository.Put(tx, apply); err != nil {
			return err
		}

		return nil
	}); err != nil {
		a.rollbackFireUserService(ctx, fireUser.UID)
		return errors.WithStack(err)
	}

	return nil
}

func (a *adminApplication) Deny(ctx context.Context, email domain.ApplyClientID) (*domain.ApplyClient, error) {
	apply, err := a.applyClientRepository.Get(ctx, email)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	apply.Deny()

	if err := a.transaction(ctx, func(tx *boom.Transaction) error {
		if err := a.applyClientRepository.Put(tx, apply); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, errors.WithStack(err)
	}

	return apply, nil
}

package customer_app

import (
	"context"

	"go.mercari.io/datastore/boom"

	"gae-go-sample/adapter"
	"gae-go-sample/domain"
	pb "gae-go-sample/proto/go/pb"

	"github.com/pkg/errors"
)

type adminApplication struct {
	executorID domain.AdminUserID
	*application
}

func (a *adminApplication) GetAll(ctx context.Context) ([]*domain.Customer, error) {
	customers, err := a.customerRepository.GetAll(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for i := range customers {
		urlWithSignature, err := a.publishResourceService(ctx, customers[i].GSIconURL)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		customers[i].SignedIconURL = urlWithSignature
	}

	return a.bindCustomerEmailService(ctx, customers)
}

func (a *adminApplication) GetAllByIDs(
	ctx context.Context,
	ids []domain.CustomerID) ([]*domain.Customer, error) {
	customers, err := a.customerRepository.GetMulti(ctx, ids)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for i := range customers {
		urlWithSignature, err := a.publishResourceService(ctx, customers[i].GSIconURL)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		customers[i].SignedIconURL = urlWithSignature
	}

	return a.bindCustomerEmailService(ctx, customers)
}

func (a *adminApplication) GetAllWithPager(ctx context.Context, page int32, offset int32) ([]*domain.Customer, error) {
	customers, err := a.customerRepository.GetAllWithPager(ctx, domain.NewPager(page, offset))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for i := range customers {
		urlWithSignature, err := a.publishResourceService(ctx, customers[i].GSIconURL)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		customers[i].SignedIconURL = urlWithSignature
	}

	return a.bindCustomerEmailService(ctx, customers)
}

func (a *adminApplication) GetAllByFilterWithPager(
	ctx context.Context,
	params adapter.FilterCustomerParams,
	page int32,
	offset int32) ([]*domain.Customer, error) {
	pager := domain.NewPager(page, offset)

	customerIDsWithHighlights, err := a.customerIndexRepository.SearchByStatusWithConditionWithPager(
		ctx,
		params.Query,
		params.Status,
		pager,
		pb.SearchCustomerRequest_OrderBy_CreatedAt_DESC)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	ids := make([]domain.CustomerID, 0, len(customerIDsWithHighlights))
	for _, id := range customerIDsWithHighlights {
		ids = append(ids, id.ID)
	}

	customers, err := a.customerRepository.GetMulti(ctx, ids)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for i := range customers {
		customerID := customers[i].ID
		for _, idWithHighlights := range customerIDsWithHighlights {
			if customerID == idWithHighlights.ID {
				customers[i].Highlights = idWithHighlights.Highlights
				break
			}
		}
	}

	for i := range customers {
		urlWithSignature, err := a.publishResourceService(ctx, customers[i].GSIconURL)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		customers[i].SignedIconURL = urlWithSignature
	}

	return a.bindCustomerEmailService(ctx, customers)
}

func (a *adminApplication) Get(
	ctx context.Context,
	id domain.CustomerID) (*domain.Customer, error) {
	customer, err := a.customerRepository.Get(ctx, id)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	fireUser, err := a.fireUserRepository.Get(ctx, domain.UserID(id))
	if err == nil {
		customer.BindEmail(fireUser)
	} else {
		a.logger.Warn().With(ctx).Printf("failed get firebase user, %#v", err)
	}

	urlWithSignature, err := a.publishResourceService(ctx, customer.GSIconURL)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	customer.SignedIconURL = urlWithSignature

	return customer, nil
}

func (a *adminApplication) GetTotalCount(ctx context.Context) (int64, error) {
	count, err := a.customerRepository.GetTotalCount(ctx)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return count, nil
}

func (a *adminApplication) GetCountByFilter(ctx context.Context, params adapter.FilterCustomerParams) (int64, error) {
	count, err := a.customerIndexRepository.SearchCountByStatusWithCondition(
		ctx,
		params.Query,
		params.Status)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return count, nil
}

func (a *adminApplication) DenyInspection(ctx context.Context, id domain.CustomerID) error {
	customer, err := a.customerRepository.Get(ctx, id)
	if err != nil {
		return errors.WithStack(err)
	}

	if err := a.transaction(ctx, func(tx *boom.Transaction) error {
		if err := a.customerRepository.Put(tx, customer); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return errors.WithStack(err)
	}

	if err := a.customerIndexRepository.Put(ctx, customer); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (a *adminApplication) PassInspection(ctx context.Context, id domain.CustomerID) error {
	customer, err := a.customerRepository.Get(ctx, id)
	if err != nil {
		return errors.WithStack(err)
	}

	if err := a.transaction(ctx, func(tx *boom.Transaction) error {
		if err := a.customerRepository.Put(tx, customer); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return errors.WithStack(err)
	}

	if err := a.customerIndexRepository.Put(ctx, customer); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (a *adminApplication) Deny(ctx context.Context, id domain.CustomerID) error {
	customer, err := a.customerRepository.Get(ctx, id)
	if err != nil {
		return errors.WithStack(err)
	}

	if err := customer.Deny(); err != nil {
		return errors.WithStack(err)
	}

	if err := a.transaction(ctx, func(tx *boom.Transaction) error {
		if err := a.customerRepository.Put(tx, customer); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return errors.WithStack(err)
	}

	if err := a.customerIndexRepository.Put(ctx, customer); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (a *adminApplication) ReInspect(ctx context.Context, id domain.CustomerID) error {
	customer, err := a.customerRepository.Get(ctx, id)
	if err != nil {
		return errors.WithStack(err)
	}

	customer.ReInspect()

	if err := a.transaction(ctx, func(tx *boom.Transaction) error {
		if err := a.customerRepository.Put(tx, customer); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return errors.WithStack(err)
	}

	if err := a.customerIndexRepository.Put(ctx, customer); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (a *adminApplication) Active(ctx context.Context, id domain.CustomerID) error {
	customer, err := a.customerRepository.Get(ctx, id)
	if err != nil {
		return errors.WithStack(err)
	}

	if err := customer.Active(); err != nil {
		return errors.WithStack(err)
	}

	if err := a.transaction(ctx, func(tx *boom.Transaction) error {
		if err := a.customerRepository.Put(tx, customer); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return errors.WithStack(err)
	}

	if err := a.customerIndexRepository.Put(ctx, customer); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

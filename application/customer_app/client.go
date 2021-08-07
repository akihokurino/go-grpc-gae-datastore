package customer_app

import (
	"context"

	"gae-go-recruiting-server/adapter"
	"gae-go-recruiting-server/domain"
	pb "gae-go-recruiting-server/proto/go/pb"

	"github.com/pkg/errors"
)

type clientApplication struct {
	executorID domain.ClientID
	*application
}

func (a *clientApplication) GetAll(ctx context.Context) ([]*domain.Customer, error) {
	customers, err := a.customerRepository.GetAll(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	customers, err = a.bindCustomerEmailService(ctx, customers)
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

	return customers, nil
}

func (a *clientApplication) GetAllByIDs(
	ctx context.Context,
	ids []domain.CustomerID) ([]*domain.Customer, error) {
	customers, err := a.customerRepository.GetMulti(ctx, ids)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	customers, err = a.bindCustomerEmailService(ctx, customers)
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

	return customers, nil
}

func (a *clientApplication) SearchWithPager(
	ctx context.Context,
	params adapter.SearchCustomerParams,
	page int32,
	offset int32) ([]*domain.Customer, error) {
	pager := domain.NewPager(page, offset)

	customerIDsWithHighlights, err := a.customerIndexRepository.SearchByStatusWithConditionWithPager(
		ctx,
		params.Query,
		pb.Customer_Status_Active,
		pager,
		params.OrderBy)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	ids := make([]domain.CustomerID, 0, len(customerIDsWithHighlights))
	for _, id := range customerIDsWithHighlights {
		ids = append(ids, id.ID)
	}

	customers, err := a.customerRepository.GetMultiWithIgnoreNotFound(ctx, ids)
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

	return customers, nil
}

func (a *clientApplication) Get(
	ctx context.Context,
	id domain.CustomerID) (*domain.Customer, error) {
	me, err := a.clientRepository.Get(ctx, a.executorID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	hasEntry, err := a.customerAlreadyEntryToAnyoneService(ctx, me, id)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if !hasEntry {
		return nil, errors.WithStack(domain.ErrCustomerDidNotEntry)
	}

	customer, err := a.customerRepository.Get(ctx, id)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if !customer.IsActive() {
		return nil, errors.WithStack(domain.ErrCustomerIsNotActive)
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

func (a *clientApplication) GetCountBySearch(
	ctx context.Context,
	params adapter.SearchCustomerParams) (int64, error) {
	count, err := a.customerIndexRepository.SearchCountByStatusWithCondition(
		ctx,
		params.Query,
		pb.Customer_Status_Active)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return count, nil
}

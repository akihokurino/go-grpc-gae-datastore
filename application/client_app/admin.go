package client_app

import (
	"context"

	"gae-go-sample/domain"

	"github.com/pkg/errors"
)

type adminApplication struct {
	executorID domain.AdminUserID
	*application
}

func (a *adminApplication) GetAllWithPager(ctx context.Context, page int32, offset int32) ([]*domain.Client, error) {
	clients, err := a.clientRepository.GetAllWithPager(ctx, domain.NewPager(page, offset))
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

func (a *adminApplication) GetAllByCompanyWithPager(
	ctx context.Context,
	companyID domain.CompanyID,
	page int32,
	offset int32) ([]*domain.Client, error) {
	clients, err := a.clientRepository.GetAllByCompanyWithPager(ctx, companyID, domain.NewPager(page, offset))
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

func (a *adminApplication) GetTotalCount(ctx context.Context) (int64, error) {
	count, err := a.clientRepository.GetTotalCount(ctx)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return count, nil
}

func (a *adminApplication) GetCountByCompany(ctx context.Context, companyID domain.CompanyID) (int64, error) {
	count, err := a.clientRepository.GetCountByCompany(ctx, companyID)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return count, nil
}

package service

import (
	"context"

	"gae-go-sample/adapter"
	"gae-go-sample/domain"

	"github.com/pkg/errors"
)

func NewValidCompanyService(companyRepository adapter.CompanyRepository) adapter.ValidCompanyService {
	return func(
		ctx context.Context,
		client *domain.Client) (*domain.Company, error) {
		company, err := companyRepository.Get(ctx, client.CompanyID)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		if !company.IsIncludeClient(client) {
			return nil, errors.WithStack(domain.ErrInvalidClient)
		}

		return company, nil
	}
}

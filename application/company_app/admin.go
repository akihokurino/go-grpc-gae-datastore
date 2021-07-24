package company_app

import (
	"context"

	"go.mercari.io/datastore/boom"

	"github.com/pkg/errors"

	"gae-go-sample/domain"
)

type adminApplication struct {
	executorID domain.AdminUserID
	*application
}

func (a *adminApplication) GetAllWithPager(ctx context.Context, page int32, offset int32) ([]*domain.Company, error) {
	companies, err := a.companyRepository.GetAllWithPager(ctx, domain.NewPager(page, offset))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for i := range companies {
		urlWithSignature, err := a.publishResourceService(ctx, companies[i].GSLogoURL)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		companies[i].SignedLogoURL = urlWithSignature
	}

	return companies, nil
}

func (a *adminApplication) GetTotalCount(ctx context.Context) (int64, error) {
	count, err := a.companyRepository.GetTotalCount(ctx)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return count, nil
}

func (a *adminApplication) Active(ctx context.Context, id domain.CompanyID) error {
	company, err := a.companyRepository.Get(ctx, id)
	if err != nil {
		return errors.WithStack(err)
	}

	company.Active()

	if err := a.transaction(ctx, func(tx *boom.Transaction) error {
		if err := a.companyRepository.Put(tx, company); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (a *adminApplication) Ban(ctx context.Context, id domain.CompanyID) error {
	company, err := a.companyRepository.Get(ctx, id)
	if err != nil {
		return errors.WithStack(err)
	}

	company.Ban()

	if err := a.transaction(ctx, func(tx *boom.Transaction) error {
		if err := a.companyRepository.Put(tx, company); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

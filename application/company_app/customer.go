package company_app

import (
	"context"

	"gae-go-sample/adapter"

	"github.com/pkg/errors"

	"gae-go-sample/domain"
)

type customerApplication struct {
	executorID domain.CustomerID
	*application
}

func (a *customerApplication) GetAllByIDsWithMaskIfNotActive(
	ctx context.Context,
	ids []domain.CompanyID) ([]*domain.Company, error) {
	companies, err := a.companyRepository.GetMulti(ctx, ids)
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

func (a *customerApplication) GetAllByIDsWithMaskIfNeedOrNotActive(
	ctx context.Context,
	ids []adapter.CompanyIDWithProjectID) ([]*adapter.CompanyWithProjectID, error) {
	companyIDs := make([]domain.CompanyID, 0, len(ids))
	for _, id := range ids {
		companyIDs = append(companyIDs, id.CompanyID)
	}

	companies, err := a.companyRepository.GetMulti(ctx, companyIDs)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	projectIDs := make([]domain.ProjectID, 0, len(ids))
	for _, id := range ids {
		projectIDs = append(projectIDs, id.ProjectID)
	}

	projects, err := a.projectRepository.GetMulti(ctx, projectIDs)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	companyMap := make(map[domain.ProjectID]*domain.Company, 0)

	for _, p := range projects {
		var company domain.Company
		for _, c := range companies {
			if c.ID == p.CompanyID {
				company = *c
			}
		}

		companyMap[p.ID] = &company
	}

	companyWithProjectIDs := make([]*adapter.CompanyWithProjectID, 0, len(companyMap))
	for key, company := range companyMap {
		companyWithProjectIDs = append(companyWithProjectIDs, &adapter.CompanyWithProjectID{
			Company:   company,
			ProjectID: key,
		})
	}

	for i := range companyWithProjectIDs {
		urlWithSignature, err := a.publishResourceService(ctx, companies[i].GSLogoURL)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		companyWithProjectIDs[i].Company.SignedLogoURL = urlWithSignature
	}

	return companyWithProjectIDs, nil
}

func (a *customerApplication) GetWithMaskIfNotActive(
	ctx context.Context,
	id domain.CompanyID) (*domain.Company, error) {
	company, err := a.companyRepository.Get(ctx, id)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	urlWithSignature, err := a.publishResourceService(ctx, company.GSLogoURL)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	company.SignedLogoURL = urlWithSignature

	return company, nil
}

func (a *customerApplication) GetWithMaskIfNeedOrNotActive(
	ctx context.Context,
	id adapter.CompanyIDWithProjectID) (*domain.Company, error) {
	company, err := a.companyRepository.Get(ctx, id.CompanyID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	urlWithSignature, err := a.publishResourceService(ctx, company.GSLogoURL)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	company.SignedLogoURL = urlWithSignature

	return company, nil
}

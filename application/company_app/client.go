package company_app

import (
	"context"

	"go.mercari.io/datastore/boom"

	"gae-go-sample/adapter"
	"gae-go-sample/domain"

	"github.com/pkg/errors"
)

type clientApplication struct {
	executorID domain.ClientID
	*application
}

func (a *clientApplication) GetAllByIDsWithMaskIfNeed(ctx context.Context, ids []adapter.CompanyIDWithProjectID) ([]*adapter.CompanyWithProjectID, error) {
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

func (a *clientApplication) GetWithMaskIfNeed(ctx context.Context, id adapter.CompanyIDWithProjectID) (*domain.Company, error) {
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

func (a *clientApplication) Update(
	ctx context.Context,
	params adapter.CompanyParams) (*domain.Company, error) {
	me, err := a.clientRepository.Get(ctx, a.executorID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	company, err := a.validCompanyService(ctx, me)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err := company.Update(
		params.Name,
		params.LogoURL,
		params.WebURL,
		params.EstablishedAt,
		params.PostalCode,
		params.RepresentativeName,
		params.CapitalStock,
		params.Introduction,
		params.AccordingCompanyName,
		params.AccordingCompanyPostalCode,
		params.AccordingCompanyAddress); err != nil {
		return nil, errors.WithStack(err)
	}

	projects, err := a.projectRepository.GetAllByCompany(ctx, company.ID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err := a.transaction(ctx, func(tx *boom.Transaction) error {
		if err := a.companyRepository.Put(tx, company); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, errors.WithStack(err)
	}

	for _, project := range projects {
		if err := a.projectIndexRepository.Put(ctx, project, company); err != nil {
			return nil, errors.WithStack(err)
		}
	}

	urlWithSignature, err := a.publishResourceService(ctx, company.GSLogoURL)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	company.SignedLogoURL = urlWithSignature

	return company, nil
}

package company_app

import (
	"context"

	"gae-go-recruiting-server/adapter"
	"gae-go-recruiting-server/domain"

	"github.com/pkg/errors"
)

type application struct {
	logger                 adapter.CompositeLogger
	companyRepository      adapter.CompanyRepository
	projectRepository      adapter.ProjectRepository
	clientRepository       adapter.ClientRepository
	customerRepository     adapter.CustomerRepository
	projectIndexRepository adapter.ProjectIndexRepository
	transaction            adapter.Transaction
	validCompanyService    adapter.ValidCompanyService
	fireUserRepository     adapter.FireUserRepository
	publishResourceService adapter.PublishResourceService
}

func NewApplication(
	logger adapter.CompositeLogger,
	companyRepository adapter.CompanyRepository,
	projectRepository adapter.ProjectRepository,
	clientRepository adapter.ClientRepository,
	customerRepository adapter.CustomerRepository,
	projectIndexRepository adapter.ProjectIndexRepository,
	transaction adapter.Transaction,
	validCompanyService adapter.ValidCompanyService,
	fireUserRepository adapter.FireUserRepository,
	publishResourceService adapter.PublishResourceService) adapter.CompanyApplication {
	return &application{
		logger:                 logger,
		companyRepository:      companyRepository,
		projectRepository:      projectRepository,
		clientRepository:       clientRepository,
		customerRepository:     customerRepository,
		projectIndexRepository: projectIndexRepository,
		transaction:            transaction,
		validCompanyService:    validCompanyService,
		fireUserRepository:     fireUserRepository,
		publishResourceService: publishResourceService,
	}
}

func (a *application) BuildAsAdmin(id domain.AdminUserID) adapter.CompanyApplicationForAdmin {
	return &adminApplication{
		executorID:  id,
		application: a,
	}
}

func (a *application) BuildAsCustomer(id domain.CustomerID) adapter.CompanyApplicationForCustomer {
	return &customerApplication{
		executorID:  id,
		application: a,
	}
}

func (a *application) BuildAsClient(id domain.ClientID) adapter.CompanyApplicationForClient {
	return &clientApplication{
		executorID:  id,
		application: a,
	}
}

func (a *application) GetAllByIDs(ctx context.Context, ids []domain.CompanyID) ([]*domain.Company, error) {
	companies, err := a.companyRepository.GetMulti(ctx, ids)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return companies, nil
}

func (a *application) Get(ctx context.Context, id domain.CompanyID) (*domain.Company, error) {
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

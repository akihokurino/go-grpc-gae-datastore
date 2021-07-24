package contract_app

import (
	"gae-go-sample/adapter"
	"gae-go-sample/domain"
)

type application struct {
	logger                 adapter.CompositeLogger
	transaction            adapter.Transaction
	contractRepository     adapter.ContractRepository
	contractFileRepository adapter.ContractFileRepository
	projectRepository      adapter.ProjectRepository
	companyRepository      adapter.CompanyRepository
	clientRepository       adapter.ClientRepository
	customerRepository     adapter.CustomerRepository
	entryRepository        adapter.EntryRepository
	fireUserRepository     adapter.FireUserRepository
	publishResourceService adapter.PublishResourceService
}

func NewApplication(
	logger adapter.CompositeLogger,
	transaction adapter.Transaction,
	contractRepository adapter.ContractRepository,
	contractFileRepository adapter.ContractFileRepository,
	projectRepository adapter.ProjectRepository,
	companyRepository adapter.CompanyRepository,
	clientRepository adapter.ClientRepository,
	customerRepository adapter.CustomerRepository,
	entryRepository adapter.EntryRepository,
	fireUserRepository adapter.FireUserRepository,
	publishResourceService adapter.PublishResourceService) adapter.ContractApplication {
	return &application{
		logger:                 logger,
		transaction:            transaction,
		contractRepository:     contractRepository,
		contractFileRepository: contractFileRepository,
		projectRepository:      projectRepository,
		companyRepository:      companyRepository,
		clientRepository:       clientRepository,
		customerRepository:     customerRepository,
		entryRepository:        entryRepository,
		fireUserRepository:     fireUserRepository,
		publishResourceService: publishResourceService,
	}
}

func (a *application) BuildAsAdmin(id domain.AdminUserID) adapter.ContractApplicationForAdmin {
	return &adminApplication{
		executorID:  id,
		application: a,
	}
}

func (a *application) BuildAsClient(id domain.ClientID) adapter.ContractApplicationForClient {
	return &clientApplication{
		executorID:  id,
		application: a,
	}
}

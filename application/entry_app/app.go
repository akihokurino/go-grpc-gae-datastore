package entry_app

import (
	"gae-go-recruiting-server/adapter"
	"gae-go-recruiting-server/domain"
)

type application struct {
	logger                      adapter.CompositeLogger
	entryRepository             adapter.EntryRepository
	projectRepository           adapter.ProjectRepository
	companyRepository           adapter.CompanyRepository
	customerRepository          adapter.CustomerRepository
	clientRepository            adapter.ClientRepository
	messageRoomRepository       adapter.MessageRoomRepository
	rtMessageRoomRepository     adapter.RtMessageRoomRepository
	noEntrySupportRepository    adapter.NoEntrySupportRepository
	noMessageSupportRepository  adapter.NoMessageSupportRepository
	transaction                 adapter.Transaction
	validProjectService         adapter.ValidProjectService
	openNoMessageSupportService adapter.OpenNoMessageSupportService
	closeNoEntrySupportService  adapter.CloseNoEntrySupportService
	fireUserRepository          adapter.FireUserRepository
}

func NewApplication(
	logger adapter.CompositeLogger,
	entryRepository adapter.EntryRepository,
	projectRepository adapter.ProjectRepository,
	companyRepository adapter.CompanyRepository,
	customerRepository adapter.CustomerRepository,
	clientRepository adapter.ClientRepository,
	messageRoomRepository adapter.MessageRoomRepository,
	rtMessageRoomRepository adapter.RtMessageRoomRepository,
	noEntrySupportRepository adapter.NoEntrySupportRepository,
	noMessageSupportRepository adapter.NoMessageSupportRepository,
	transaction adapter.Transaction,
	validProjectService adapter.ValidProjectService,
	openNoMessageSupportService adapter.OpenNoMessageSupportService,
	closeNoEntrySupportService adapter.CloseNoEntrySupportService,
	fireUserRepository adapter.FireUserRepository) adapter.EntryApplication {
	return &application{
		logger:                      logger,
		entryRepository:             entryRepository,
		projectRepository:           projectRepository,
		companyRepository:           companyRepository,
		customerRepository:          customerRepository,
		clientRepository:            clientRepository,
		messageRoomRepository:       messageRoomRepository,
		rtMessageRoomRepository:     rtMessageRoomRepository,
		noEntrySupportRepository:    noEntrySupportRepository,
		noMessageSupportRepository:  noMessageSupportRepository,
		transaction:                 transaction,
		validProjectService:         validProjectService,
		openNoMessageSupportService: openNoMessageSupportService,
		closeNoEntrySupportService:  closeNoEntrySupportService,
		fireUserRepository:          fireUserRepository,
	}
}

func (a *application) BuildAsAdmin(id domain.AdminUserID) adapter.EntryApplicationForAdmin {
	return &adminApplication{
		executorID:  id,
		application: a,
	}
}

func (a *application) BuildAsCustomer(id domain.CustomerID) adapter.EntryApplicationForCustomer {
	return &customerApplication{
		executorID:  id,
		application: a,
	}
}

func (a *application) BuildAsClient(id domain.ClientID) adapter.EntryApplicationForClient {
	return &clientApplication{
		executorID:  id,
		application: a,
	}
}

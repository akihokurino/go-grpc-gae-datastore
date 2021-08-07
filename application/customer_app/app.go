package customer_app

import (
	"gae-go-recruiting-server/adapter"
	"gae-go-recruiting-server/domain"
)

type application struct {
	logger                              adapter.CompositeLogger
	userRepository                      adapter.UserRepository
	customerRepository                  adapter.CustomerRepository
	customerIndexRepository             adapter.CustomerIndexRepository
	clientRepository                    adapter.ClientRepository
	companyRepository                   adapter.CompanyRepository
	fireUserRepository                  adapter.FireUserRepository
	projectRepository                   adapter.ProjectRepository
	entryRepository                     adapter.EntryRepository
	messageRoomRepository               adapter.MessageRoomRepository
	rtMemberRepository                  adapter.RtMemberRepository
	transaction                         adapter.Transaction
	bindCustomerEmailService            adapter.BindCustomerEmailService
	customerAlreadyEntryToAnyoneService adapter.CustomerAlreadyEntryToAnyoneService
	rollbackFireUserService             adapter.RollbackFireUserService
	publishResourceService              adapter.PublishResourceService
}

func NewApplication(
	logger adapter.CompositeLogger,
	userRepository adapter.UserRepository,
	customerRepository adapter.CustomerRepository,
	customerIndexRepository adapter.CustomerIndexRepository,
	clientRepository adapter.ClientRepository,
	companyRepository adapter.CompanyRepository,
	fireUserRepository adapter.FireUserRepository,
	projectRepository adapter.ProjectRepository,
	entryRepository adapter.EntryRepository,
	messageRoomRepository adapter.MessageRoomRepository,
	rtMemberRepository adapter.RtMemberRepository,
	transaction adapter.Transaction,
	bindCustomerEmailService adapter.BindCustomerEmailService,
	customerAlreadyEntryToAnyoneService adapter.CustomerAlreadyEntryToAnyoneService,
	rollbackFireUserService adapter.RollbackFireUserService,
	publishResourceService adapter.PublishResourceService) adapter.CustomerApplication {
	return &application{
		logger:                              logger,
		userRepository:                      userRepository,
		customerRepository:                  customerRepository,
		customerIndexRepository:             customerIndexRepository,
		clientRepository:                    clientRepository,
		companyRepository:                   companyRepository,
		fireUserRepository:                  fireUserRepository,
		projectRepository:                   projectRepository,
		entryRepository:                     entryRepository,
		messageRoomRepository:               messageRoomRepository,
		rtMemberRepository:                  rtMemberRepository,
		transaction:                         transaction,
		bindCustomerEmailService:            bindCustomerEmailService,
		customerAlreadyEntryToAnyoneService: customerAlreadyEntryToAnyoneService,
		rollbackFireUserService:             rollbackFireUserService,
		publishResourceService:              publishResourceService,
	}
}

func (a *application) BuildAsAdmin(id domain.AdminUserID) adapter.CustomerApplicationForAdmin {
	return &adminApplication{
		executorID:  id,
		application: a,
	}
}

func (a *application) BuildAsCustomer(id domain.CustomerID) adapter.CustomerApplicationForCustomer {
	return &customerApplication{
		executorID:  id,
		application: a,
	}
}

func (a *application) BuildAsClient(id domain.ClientID) adapter.CustomerApplicationForClient {
	return &clientApplication{
		executorID:  id,
		application: a,
	}
}

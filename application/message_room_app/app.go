package message_room_app

import (
	"gae-go-recruiting-server/adapter"
	"gae-go-recruiting-server/domain"
)

type application struct {
	logger                     adapter.CompositeLogger
	userRepository             adapter.UserRepository
	companyRepository          adapter.CompanyRepository
	customerRepository         adapter.CustomerRepository
	clientRepository           adapter.ClientRepository
	projectRepository          adapter.ProjectRepository
	messageRoomRepository      adapter.MessageRoomRepository
	rtMessageRoomRepository    adapter.RtMessageRoomRepository
	messageRepository          adapter.MessageRepository
	noMessageSupportRepository adapter.NoMessageSupportRepository
	transaction                adapter.Transaction
}

func NewApplication(
	logger adapter.CompositeLogger,
	userRepository adapter.UserRepository,
	companyRepository adapter.CompanyRepository,
	customerRepository adapter.CustomerRepository,
	clientRepository adapter.ClientRepository,
	projectRepository adapter.ProjectRepository,
	messageRoomRepository adapter.MessageRoomRepository,
	rtMessageRoomRepository adapter.RtMessageRoomRepository,
	messageRepository adapter.MessageRepository,
	noMessageSupportRepository adapter.NoMessageSupportRepository,
	transaction adapter.Transaction) adapter.MessageRoomApplication {
	return &application{
		logger:                     logger,
		userRepository:             userRepository,
		companyRepository:          companyRepository,
		customerRepository:         customerRepository,
		clientRepository:           clientRepository,
		projectRepository:          projectRepository,
		messageRoomRepository:      messageRoomRepository,
		rtMessageRoomRepository:    rtMessageRoomRepository,
		messageRepository:          messageRepository,
		noMessageSupportRepository: noMessageSupportRepository,
		transaction:                transaction,
	}
}

func (a *application) BuildAsAdmin(id domain.AdminUserID) adapter.MessageRoomApplicationForAdmin {
	return &adminApplication{
		executorID:  id,
		application: a,
	}
}

func (a *application) BuildAsCustomer(id domain.CustomerID) adapter.MessageRoomApplicationForCustomer {
	return &customerApplication{
		executorID:  id,
		application: a,
	}
}

func (a *application) BuildAsClient(id domain.ClientID) adapter.MessageRoomApplicationForClient {
	return &clientApplication{
		executorID:  id,
		application: a,
	}
}

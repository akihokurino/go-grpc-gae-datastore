package message_app

import (
	"context"
	"time"

	"go.mercari.io/datastore/boom"

	"gae-go-sample/adapter"
	"gae-go-sample/domain"

	"github.com/pkg/errors"
)

type application struct {
	logger                       adapter.CompositeLogger
	userRepository               adapter.UserRepository
	messageRoomRepository        adapter.MessageRoomRepository
	messageRepository            adapter.MessageRepository
	customerRepository           adapter.CustomerRepository
	projectRepository            adapter.ProjectRepository
	noMessageSupportRepository   adapter.NoMessageSupportRepository
	companyRepository            adapter.CompanyRepository
	clientRepository             adapter.ClientRepository
	transaction                  adapter.Transaction
	closeNoMessageSupportService adapter.CloseNoMessageSupportService
	fireUserRepository           adapter.FireUserRepository
	publishResourceService       adapter.PublishResourceService
}

func NewApplication(
	logger adapter.CompositeLogger,
	userRepository adapter.UserRepository,
	messageRoomRepository adapter.MessageRoomRepository,
	messageRepository adapter.MessageRepository,
	customerRepository adapter.CustomerRepository,
	projectRepository adapter.ProjectRepository,
	noMessageSupportRepository adapter.NoMessageSupportRepository,
	companyRepository adapter.CompanyRepository,
	clientRepository adapter.ClientRepository,
	transaction adapter.Transaction,
	closeNoMessageSupportService adapter.CloseNoMessageSupportService,
	fireUserRepository adapter.FireUserRepository,
	publishResourceService adapter.PublishResourceService) adapter.MessageApplication {
	return &application{
		logger:                       logger,
		userRepository:               userRepository,
		messageRoomRepository:        messageRoomRepository,
		messageRepository:            messageRepository,
		customerRepository:           customerRepository,
		projectRepository:            projectRepository,
		noMessageSupportRepository:   noMessageSupportRepository,
		companyRepository:            companyRepository,
		clientRepository:             clientRepository,
		transaction:                  transaction,
		closeNoMessageSupportService: closeNoMessageSupportService,
		fireUserRepository:           fireUserRepository,
		publishResourceService:       publishResourceService,
	}
}

func (a *application) BuildAsAdmin(id domain.AdminUserID) adapter.MessageApplicationForAdmin {
	return &adminApplication{
		executorID:  id,
		application: a,
	}
}

func (a *application) BuildAsCustomer(id domain.CustomerID) adapter.MessageApplicationForCustomer {
	return &customerApplication{
		executorID:  id,
		application: a,
	}
}

func (a *application) BuildAsClient(id domain.ClientID) adapter.MessageApplicationForClient {
	return &clientApplication{
		executorID:  id,
		application: a,
	}
}

func (a *application) Create(
	ctx context.Context,
	messageID domain.MessageID,
	roomID domain.MessageRoomID,
	fromID string,
	toID string,
	params adapter.MessageParams,
	now time.Time) error {
	room, err := a.messageRoomRepository.Get(ctx, roomID)
	if err != nil {
		return errors.WithStack(err)
	}

	var to domain.MessageRoomUser
	fromCompany := false

	_, err = a.userRepository.Get(ctx, domain.UserID(fromID))
	if err != nil && !domain.IsNoSuchEntityErr(err) {
		return errors.WithStack(err)
	}
	if err != nil && domain.IsNoSuchEntityErr(err) {
		_, err := a.companyRepository.Get(ctx, domain.CompanyID(fromID))
		if err != nil {
			return errors.WithStack(err)
		}

		to = domain.ByCustomer
		fromCompany = true
	} else {
		to = domain.ByCompany
	}

	message, err := domain.NewMessage(
		messageID,
		room.ID(),
		fromID,
		toID,
		fromCompany,
		params.Text,
		params.ImageURL,
		params.FileURL,
		now)
	if err != nil {
		return errors.WithStack(err)
	}

	if err := a.transaction(ctx, func(tx *boom.Transaction) error {
		if err := a.messageRepository.Put(tx, message); err != nil {
			return err
		}

		if err := a.closeNoMessageSupportService(
			ctx,
			tx,
			room); err != nil {
			return err
		}

		if err := room.ReceiveMessage(to, now); err != nil {
			return err
		}

		if err := a.messageRoomRepository.Put(tx, room); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

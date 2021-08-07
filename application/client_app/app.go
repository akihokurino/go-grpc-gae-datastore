package client_app

import (
	"context"

	"gae-go-recruiting-server/adapter"
	"gae-go-recruiting-server/domain"

	"github.com/pkg/errors"
)

type application struct {
	logger                  adapter.CompositeLogger
	userRepository          adapter.UserRepository
	clientRepository        adapter.ClientRepository
	fireUserRepository      adapter.FireUserRepository
	messageRoomRepository   adapter.MessageRoomRepository
	rtMemberRepository      adapter.RtMemberRepository
	transaction             adapter.Transaction
	bindClientEmail         adapter.BindClientEmailService
	rollbackFireUserService adapter.RollbackFireUserService
	companyRepository       adapter.CompanyRepository
	publishResourceService  adapter.PublishResourceService
}

func NewApplication(
	logger adapter.CompositeLogger,
	userRepository adapter.UserRepository,
	clientRepository adapter.ClientRepository,
	fireUserRepository adapter.FireUserRepository,
	messageRoomRepository adapter.MessageRoomRepository,
	rtMemberRepository adapter.RtMemberRepository,
	transaction adapter.Transaction,
	bindClientEmail adapter.BindClientEmailService,
	rollbackFireUserService adapter.RollbackFireUserService,
	companyRepository adapter.CompanyRepository,
	publishResourceService adapter.PublishResourceService) adapter.ClientApplication {
	return &application{
		logger:                  logger,
		userRepository:          userRepository,
		clientRepository:        clientRepository,
		fireUserRepository:      fireUserRepository,
		messageRoomRepository:   messageRoomRepository,
		rtMemberRepository:      rtMemberRepository,
		transaction:             transaction,
		bindClientEmail:         bindClientEmail,
		rollbackFireUserService: rollbackFireUserService,
		companyRepository:       companyRepository,
		publishResourceService:  publishResourceService,
	}
}

func (a *application) BuildAsAdmin(id domain.AdminUserID) adapter.ClientApplicationForAdmin {
	return &adminApplication{
		executorID:  id,
		application: a,
	}
}

func (a *application) BuildAsClient(id domain.ClientID) adapter.ClientApplicationForClient {
	return &clientApplication{
		executorID:  id,
		application: a,
	}
}

func (a *application) Get(ctx context.Context, id domain.ClientID) (*domain.Client, error) {
	client, err := a.clientRepository.Get(ctx, id)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	fireUser, err := a.fireUserRepository.Get(ctx, domain.UserID(client.ID))
	if err == nil {
		client.BindEmail(fireUser)
	} else {
		a.logger.Error().With(ctx).Printf("failed get firebase user, %#v", err)
	}

	urlWithSignature, err := a.publishResourceService(ctx, client.GSIconURL)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	client.SignedIconURL = urlWithSignature

	return client, nil
}

package user_app

import (
	"context"

	"gae-go-sample/adapter"
	"gae-go-sample/domain"
	pb "gae-go-sample/proto/go/pb"

	"github.com/pkg/errors"
)

type application struct {
	logger                 adapter.CompositeLogger
	userRepository         adapter.UserRepository
	customerRepository     adapter.CustomerRepository
	clientRepository       adapter.ClientRepository
	fireUserRepository     adapter.FireUserRepository
	messageRoomRepository  adapter.MessageRoomRepository
	publishResourceService adapter.PublishResourceService
}

func NewApplication(
	logger adapter.CompositeLogger,
	userRepository adapter.UserRepository,
	customerRepository adapter.CustomerRepository,
	clientRepository adapter.ClientRepository,
	fireUserRepository adapter.FireUserRepository,
	messageRoomRepository adapter.MessageRoomRepository,
	publishResourceService adapter.PublishResourceService) adapter.UserApplication {
	return &application{
		logger:                 logger,
		userRepository:         userRepository,
		customerRepository:     customerRepository,
		clientRepository:       clientRepository,
		fireUserRepository:     fireUserRepository,
		messageRoomRepository:  messageRoomRepository,
		publishResourceService: publishResourceService,
	}
}

func (a *application) Get(ctx context.Context, id domain.UserID) (*domain.Me, error) {
	user, err := a.userRepository.Get(ctx, id)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	fireUser, err := a.fireUserRepository.Get(ctx, user.ID)
	if err == nil {
		user.BindEmail(fireUser)
	} else {
		a.logger.Warn().With(ctx).Printf("failed get firebase user, %#v", err)
	}

	a.logger.Info().With(ctx).Printf("current user is %#v", user)

	switch user.Role {
	case pb.User_Role_Customer:
		customer, err := a.customerRepository.Get(ctx, user.CustomerID())
		if err != nil {
			return nil, errors.WithStack(err)
		}

		rooms, err := a.messageRoomRepository.GetAllByCustomer(ctx, customer.ID)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		me := domain.NewMe(user, customer, nil)

		if err := me.DetectUnRead(rooms); err != nil {
			return nil, errors.WithStack(err)
		}

		urlWithSignature, err := a.publishResourceService(ctx, me.Customer.GSIconURL)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		me.Customer.SignedIconURL = urlWithSignature

		a.logger.Info().With(ctx).Printf("current customer is %#v", me)

		return me, nil
	case pb.User_Role_Client:
		client, err := a.clientRepository.Get(ctx, user.ClientID())
		if err != nil {
			return nil, errors.WithStack(err)
		}

		rooms, err := a.messageRoomRepository.GetAllByCompany(ctx, client.CompanyID)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		me := domain.NewMe(user, nil, client)

		if err := me.DetectUnRead(rooms); err != nil {
			return nil, errors.WithStack(err)
		}

		urlWithSignature, err := a.publishResourceService(ctx, me.Client.GSIconURL)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		me.Client.SignedIconURL = urlWithSignature

		a.logger.Info().With(ctx).Printf("current client is %#v", me)

		return me, nil
	default:
		return nil, errors.WithStack(domain.ErrInvalidUserRole)
	}
}

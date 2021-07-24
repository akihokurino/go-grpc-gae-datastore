package api

import (
	"context"

	"gae-go-sample/adapter"
	"gae-go-sample/domain"
	"gae-go-sample/handler"
	pb "gae-go-sample/proto/go/pb"
)

type messageRoomHandler struct {
	errorConverter         adapter.ErrorConverter
	messageRoomApplication adapter.MessageRoomApplication
	userApplication        adapter.UserApplication
	contextProvider        handler.ContextProvider
}

func NewMessageRoomHandler(
	errorConverter adapter.ErrorConverter,
	messageRoomApplication adapter.MessageRoomApplication,
	userApplication adapter.UserApplication,
	contextProvider handler.ContextProvider) pb.MessageRoomService {
	return &messageRoomHandler{
		errorConverter:         errorConverter,
		messageRoomApplication: messageRoomApplication,
		userApplication:        userApplication,
		contextProvider:        contextProvider,
	}
}

func (h *messageRoomHandler) GetAllByCustomer(ctx context.Context, req *pb.Pager) (*pb.MessageRoomList, error) {
	customerID := domain.CustomerID(h.contextProvider.MustAuthUID(ctx))
	app := h.messageRoomApplication.BuildAsCustomer(customerID)

	rooms, err := app.GetAllWithPager(
		ctx,
		req.Page,
		req.Offset)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	responses := make([]*pb.MessageRoom, 0, len(rooms))
	for _, room := range rooms {
		responses = append(responses, handler.ToMessageRoomResponse(room))
	}

	return &pb.MessageRoomList{
		Items: responses,
	}, nil
}

func (h *messageRoomHandler) GetAllByClient(ctx context.Context, req *pb.Pager) (*pb.MessageRoomList, error) {
	clientID := domain.ClientID(h.contextProvider.MustAuthUID(ctx))
	app := h.messageRoomApplication.BuildAsClient(clientID)

	rooms, err := app.GetAllWithPager(
		ctx,
		req.Page,
		req.Offset)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	responses := make([]*pb.MessageRoom, 0, len(rooms))
	for _, room := range rooms {
		responses = append(responses, handler.ToMessageRoomResponse(room))
	}

	return &pb.MessageRoomList{
		Items: responses,
	}, nil
}

func (h *messageRoomHandler) Get(ctx context.Context, req *pb.GetMessageRoomRequest) (*pb.MessageRoom, error) {
	userID := h.contextProvider.MustAuthUID(ctx)

	user, err := h.userApplication.Get(ctx, userID)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	var room *domain.MessageRoom

	switch user.Role {
	case pb.User_Role_Customer:
		app := h.messageRoomApplication.BuildAsCustomer(user.CustomerID())

		room, err = app.Get(
			ctx,
			domain.ProjectID(req.ProjectID),
			domain.CompanyID(req.OpponentID))
	case pb.User_Role_Client:
		app := h.messageRoomApplication.BuildAsClient(user.ClientID())

		room, err = app.Get(
			ctx,
			domain.ProjectID(req.ProjectID),
			domain.CustomerID(req.OpponentID))
	default:
		return nil, h.errorConverter(ctx, domain.ErrInvalidUserRole)
	}

	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return handler.ToMessageRoomResponse(room), nil
}

func (h *messageRoomHandler) Read(ctx context.Context, req *pb.ReadMessageRoomRequest) (*pb.Empty, error) {
	userID := h.contextProvider.MustAuthUID(ctx)

	user, err := h.userApplication.Get(ctx, userID)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	switch user.Role {
	case pb.User_Role_Customer:
		app := h.messageRoomApplication.BuildAsCustomer(user.CustomerID())

		err = app.Read(
			ctx,
			domain.ProjectID(req.ProjectID),
			domain.CompanyID(req.OpponentID))
	case pb.User_Role_Client:
		app := h.messageRoomApplication.BuildAsClient(user.ClientID())

		err = app.Read(
			ctx,
			domain.ProjectID(req.ProjectID),
			domain.CustomerID(req.OpponentID))
	default:
		return nil, h.errorConverter(ctx, domain.ErrInvalidUserRole)
	}

	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return &pb.Empty{}, nil
}

func (h *messageRoomHandler) Delete(ctx context.Context, req *pb.DeleteMessageRoomRequest) (*pb.Empty, error) {
	userID := h.contextProvider.MustAuthUID(ctx)

	user, err := h.userApplication.Get(ctx, userID)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	switch user.Role {
	case pb.User_Role_Customer:
		app := h.messageRoomApplication.BuildAsCustomer(user.CustomerID())

		err = app.Delete(
			ctx,
			domain.ProjectID(req.ProjectID),
			domain.CompanyID(req.OpponentID))
	case pb.User_Role_Client:
		app := h.messageRoomApplication.BuildAsClient(user.ClientID())

		err = app.Delete(
			ctx,
			domain.ProjectID(req.ProjectID),
			domain.CustomerID(req.OpponentID))
	default:
		return nil, h.errorConverter(ctx, domain.ErrInvalidUserRole)
	}

	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return &pb.Empty{}, nil
}

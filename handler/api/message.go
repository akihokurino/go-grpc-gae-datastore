package api

import (
	"context"

	"gae-go-recruiting-server/adapter"
	"gae-go-recruiting-server/domain"
	"gae-go-recruiting-server/handler"
	pb "gae-go-recruiting-server/proto/go/pb"
)

type messageHandler struct {
	errorConverter     adapter.ErrorConverter
	messageApplication adapter.MessageApplication
	userApplication    adapter.UserApplication
	contextProvider    handler.ContextProvider
}

func NewMessageHandler(
	errorConverter adapter.ErrorConverter,
	messageApplication adapter.MessageApplication,
	userApplication adapter.UserApplication,
	contextProvider handler.ContextProvider) pb.MessageService {
	return &messageHandler{
		errorConverter:     errorConverter,
		messageApplication: messageApplication,
		userApplication:    userApplication,
		contextProvider:    contextProvider,
	}
}

func (h *messageHandler) GetAllByRoom(ctx context.Context, req *pb.GetAllMessageByRoomRequest) (*pb.MessageList, error) {
	userID := h.contextProvider.MustAuthUID(ctx)

	user, err := h.userApplication.Get(ctx, userID)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	var messages []*domain.Message

	switch user.Role {
	case pb.User_Role_Customer:
		app := h.messageApplication.BuildAsCustomer(user.CustomerID())

		messages, err = app.GetAllByRoomWithPager(
			ctx,
			domain.ProjectID(req.ProjectID),
			domain.CompanyID(req.OpponentID),
			req.Pager.Page,
			req.Pager.Offset)
	case pb.User_Role_Client:
		app := h.messageApplication.BuildAsClient(user.ClientID())

		messages, err = app.GetAllByRoomWithPager(
			ctx,
			domain.ProjectID(req.ProjectID),
			domain.CustomerID(req.OpponentID),
			req.Pager.Page,
			req.Pager.Offset)
	default:
		return nil, h.errorConverter(ctx, domain.ErrInvalidUserRole)
	}

	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	responses := make([]*pb.Message, 0, len(messages))
	for _, message := range messages {
		responses = append(responses, handler.ToMessageResponse(message))
	}

	return &pb.MessageList{
		Items: responses,
	}, nil
}

func (h *messageHandler) GetAllNewestByRooms(ctx context.Context, req *pb.MessageRoomPartialIDList) (*pb.MessageList, error) {
	userID := h.contextProvider.MustAuthUID(ctx)

	user, err := h.userApplication.Get(ctx, userID)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	var messages []*domain.Message

	switch user.Role {
	case pb.User_Role_Customer:
		app := h.messageApplication.BuildAsCustomer(user.CustomerID())

		messageIDParams := make([]struct {
			ProjectID domain.ProjectID
			CompanyID domain.CompanyID
		}, 0, len(req.Ids))
		for _, id := range req.Ids {
			messageIDParams = append(messageIDParams, struct {
				ProjectID domain.ProjectID
				CompanyID domain.CompanyID
			}{
				ProjectID: domain.ProjectID(id.ProjectID),
				CompanyID: domain.CompanyID(id.OpponentID),
			})
		}

		messages, err = app.GetAllNewestByRooms(ctx, messageIDParams)
	case pb.User_Role_Client:
		app := h.messageApplication.BuildAsClient(user.ClientID())

		messageIDParams := make([]struct {
			ProjectID  domain.ProjectID
			CustomerID domain.CustomerID
		}, 0, len(req.Ids))
		for _, id := range req.Ids {
			messageIDParams = append(messageIDParams, struct {
				ProjectID  domain.ProjectID
				CustomerID domain.CustomerID
			}{
				ProjectID:  domain.ProjectID(id.ProjectID),
				CustomerID: domain.CustomerID(id.OpponentID),
			})
		}

		messages, err = app.GetAllNewestByRooms(ctx, messageIDParams)
	default:
		return nil, h.errorConverter(ctx, domain.ErrInvalidUserRole)
	}

	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	responses := make([]*pb.Message, 0, len(messages))
	for _, message := range messages {
		if message == nil {
			continue
		}

		responses = append(responses, handler.ToMessageResponse(message))
	}

	return &pb.MessageList{
		Items: responses,
	}, nil
}

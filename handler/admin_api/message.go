package admin_api

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
	contextProvider    handler.ContextProvider
}

func NewMessageHandler(
	errorConverter adapter.ErrorConverter,
	messageApplication adapter.MessageApplication,
	contextProvider handler.ContextProvider) pb.AdminMessageService {
	return &messageHandler{
		errorConverter:     errorConverter,
		messageApplication: messageApplication,
		contextProvider:    contextProvider,
	}
}

func (h *messageHandler) GetAllByRoom(ctx context.Context, req *pb.MessageRoomIDWithPager) (*pb.MessageList, error) {
	username := h.contextProvider.MustAuthAdminUID(ctx)
	app := h.messageApplication.BuildAsAdmin(username)

	roomID := domain.MessageRoomID(req.Id)

	messages, err := app.GetAllByRoomWithPager(
		ctx,
		roomID,
		req.Pager.Page,
		req.Pager.Offset)
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

func (h *messageHandler) GetAllNewestByRooms(ctx context.Context, req *pb.MessageRoomIDList) (*pb.MessageList, error) {
	username := h.contextProvider.MustAuthAdminUID(ctx)
	app := h.messageApplication.BuildAsAdmin(username)

	roomIDs := make([]domain.MessageRoomID, 0, len(req.Ids))
	for _, id := range req.Ids {
		roomIDs = append(roomIDs, domain.MessageRoomID(id))
	}

	messages, err := app.GetAllNewestByRooms(ctx, roomIDs)
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

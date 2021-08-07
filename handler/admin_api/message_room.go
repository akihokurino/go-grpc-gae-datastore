package admin_api

import (
	"context"

	"gae-go-recruiting-server/adapter"
	"gae-go-recruiting-server/domain"
	"gae-go-recruiting-server/handler"
	pb "gae-go-recruiting-server/proto/go/pb"
)

type messageRoomHandler struct {
	errorConverter         adapter.ErrorConverter
	messageRoomApplication adapter.MessageRoomApplication
	contextProvider        handler.ContextProvider
}

func NewMessageRoomHandler(
	errorConverter adapter.ErrorConverter,
	messageRoomApplication adapter.MessageRoomApplication,
	contextProvider handler.ContextProvider) pb.AdminMessageRoomService {
	return &messageRoomHandler{
		errorConverter:         errorConverter,
		messageRoomApplication: messageRoomApplication,
		contextProvider:        contextProvider,
	}
}

func (h *messageRoomHandler) GetAllByIDs(ctx context.Context, req *pb.MessageRoomIDWithoutClientIDList) (*pb.MessageRoomList, error) {
	username := h.contextProvider.MustAuthAdminUID(ctx)
	app := h.messageRoomApplication.BuildAsAdmin(username)

	roomIDs := make([]struct {
		ProjectID  domain.ProjectID
		CustomerID domain.CustomerID
	}, 0, len(req.Ids))
	for _, id := range req.Ids {
		roomIDs = append(roomIDs, struct {
			ProjectID  domain.ProjectID
			CustomerID domain.CustomerID
		}{
			ProjectID:  domain.ProjectID(id.ProjectID),
			CustomerID: domain.CustomerID(id.CustomerID),
		})
	}

	rooms, err := app.GetAllByIDs(ctx, roomIDs)
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

package api

import (
	"context"

	"gae-go-recruiting-server/adapter"
	"gae-go-recruiting-server/domain"
	"gae-go-recruiting-server/handler"
	pb "gae-go-recruiting-server/proto/go/pb"
)

type userHandler struct {
	errorConverter  adapter.ErrorConverter
	userApplication adapter.UserApplication
	contextProvider handler.ContextProvider
}

func NewUserHandler(
	errorConverter adapter.ErrorConverter,
	userApplication adapter.UserApplication,
	contextProvider handler.ContextProvider) pb.UserService {
	return &userHandler{
		errorConverter:  errorConverter,
		userApplication: userApplication,
		contextProvider: contextProvider,
	}
}

func (h *userHandler) GetMe(ctx context.Context, req *pb.Empty) (*pb.Me, error) {
	userID := h.contextProvider.MustAuthUID(ctx)

	user, err := h.userApplication.Get(ctx, userID)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	stateList := make([]*pb.Me_MessageState, 0, len(user.MessageStateList))
	for _, state := range user.MessageStateList {
		stateList = append(stateList, &pb.Me_MessageState{
			RoomID:      state.RoomID.String(),
			IsUnRead:    state.IsUnRead,
			UnReadCount: state.UnReadCount,
		})
	}

	switch user.Role {
	case pb.User_Role_Customer:
		return &pb.Me{
			Email:            user.Email,
			Role:             user.Role,
			Customer:         handler.ToCustomerResponse(user.Customer),
			MessageStateList: stateList,
		}, nil
	case pb.User_Role_Client:
		return &pb.Me{
			Email:            user.Email,
			Role:             user.Role,
			Client:           handler.ToClientResponse(user.Client),
			MessageStateList: stateList,
		}, nil
	default:
		return nil, h.errorConverter(ctx, domain.ErrInvalidUserRole)
	}
}

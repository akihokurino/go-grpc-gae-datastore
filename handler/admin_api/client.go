package admin_api

import (
	"context"

	"gae-go-recruiting-server/adapter"
	"gae-go-recruiting-server/domain"
	"gae-go-recruiting-server/handler"
	pb "gae-go-recruiting-server/proto/go/pb"
)

type clientHandler struct {
	errorConverter    adapter.ErrorConverter
	clientApplication adapter.ClientApplication
	contextProvider   handler.ContextProvider
}

func NewClientHandler(
	errorConverter adapter.ErrorConverter,
	clientApplication adapter.ClientApplication,
	contextProvider handler.ContextProvider) pb.AdminClientService {
	return &clientHandler{
		errorConverter:    errorConverter,
		clientApplication: clientApplication,
		contextProvider:   contextProvider,
	}
}

func (h *clientHandler) GetAll(ctx context.Context, req *pb.Pager) (*pb.ClientList, error) {
	username := h.contextProvider.MustAuthAdminUID(ctx)
	app := h.clientApplication.BuildAsAdmin(username)

	clients, err := app.GetAllWithPager(ctx, req.Page, req.Offset)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	responses := make([]*pb.Client, 0, len(clients))
	for _, client := range clients {
		responses = append(responses, handler.ToClientResponse(client))
	}

	return &pb.ClientList{
		Items: responses,
	}, nil
}

func (h *clientHandler) GetAllByCompany(ctx context.Context, req *pb.CompanyIDWithPager) (*pb.ClientList, error) {
	username := h.contextProvider.MustAuthAdminUID(ctx)
	app := h.clientApplication.BuildAsAdmin(username)

	clients, err := app.GetAllByCompanyWithPager(ctx, domain.CompanyID(req.Id), req.Pager.Page, req.Pager.Offset)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	responses := make([]*pb.Client, 0, len(clients))
	for _, client := range clients {
		responses = append(responses, handler.ToClientResponse(client))
	}

	return &pb.ClientList{
		Items: responses,
	}, nil
}

func (h *clientHandler) Get(ctx context.Context, req *pb.ClientID) (*pb.Client, error) {
	username := h.contextProvider.MustAuthAdminUID(ctx)
	app := h.clientApplication.BuildAsAdmin(username)

	client, err := app.Get(ctx, domain.ClientID(req.Id))
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return handler.ToClientResponse(client), nil
}

func (h *clientHandler) GetTotalCount(ctx context.Context, req *pb.Empty) (*pb.Count, error) {
	username := h.contextProvider.MustAuthAdminUID(ctx)
	app := h.clientApplication.BuildAsAdmin(username)

	count, err := app.GetTotalCount(ctx)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return &pb.Count{Count: count}, nil
}

func (h *clientHandler) GetCountByCompany(ctx context.Context, req *pb.CompanyID) (*pb.Count, error) {
	username := h.contextProvider.MustAuthAdminUID(ctx)
	app := h.clientApplication.BuildAsAdmin(username)

	count, err := app.GetCountByCompany(ctx, domain.CompanyID(req.Id))
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return &pb.Count{Count: count}, nil
}

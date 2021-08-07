package api

import (
	"context"
	"net/url"
	"time"

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
	contextProvider handler.ContextProvider) pb.ClientService {
	return &clientHandler{
		errorConverter:    errorConverter,
		clientApplication: clientApplication,
		contextProvider:   contextProvider,
	}
}

func (h *clientHandler) GetAll(ctx context.Context, req *pb.Pager) (*pb.ClientList, error) {
	clientID := domain.ClientID(h.contextProvider.MustAuthUID(ctx))
	app := h.clientApplication.BuildAsClient(clientID)

	clients, err := app.GetAll(ctx, req.Page, req.Offset)
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

func (h *clientHandler) GetAllByIDs(ctx context.Context, req *pb.ClientIDList) (*pb.ClientList, error) {
	clientID := domain.ClientID(h.contextProvider.MustAuthUID(ctx))
	app := h.clientApplication.BuildAsClient(clientID)

	ids := make([]domain.ClientID, 0, len(req.Ids))
	for _, id := range req.Ids {
		ids = append(ids, domain.ClientID(id))
	}

	clients, err := app.GetAllByIDs(ctx, ids)
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
	clientID := domain.ClientID(h.contextProvider.MustAuthUID(ctx))
	app := h.clientApplication.BuildAsClient(clientID)

	client, err := app.Get(ctx, domain.ClientID(req.Id))
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return handler.ToClientResponse(client), nil
}

func (h *clientHandler) Create(ctx context.Context, req *pb.CreateClientRequest) (*pb.Client, error) {
	clientID := domain.ClientID(h.contextProvider.MustAuthUID(ctx))
	app := h.clientApplication.BuildAsClient(clientID)

	now := time.Now()

	client, err := app.Create(
		ctx,
		req.Email,
		req.Password,
		adapter.CreateClientParams{
			Name:        req.Name,
			NameKana:    req.NameKana,
			PhoneNumber: req.PhoneNumber,
			Position:    req.Position,
			Role:        req.Role,
		},
		now,
	)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return handler.ToClientResponse(client), nil
}

func (h *clientHandler) Update(ctx context.Context, req *pb.UpdateClientRequest) (*pb.Client, error) {
	clientID := domain.ClientID(h.contextProvider.MustAuthUID(ctx))
	app := h.clientApplication.BuildAsClient(clientID)

	iconURL, _ := url.Parse(req.IconURL)

	client, err := app.Update(
		ctx,
		adapter.UpdateClientParams{
			Name:        req.Name,
			NameKana:    req.NameKana,
			IconURL:     iconURL,
			PhoneNumber: req.PhoneNumber,
			Position:    req.Position,
		},
	)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return handler.ToClientResponse(client), nil
}

func (h *clientHandler) UpdateRole(ctx context.Context, req *pb.UpdateClientRoleRequest) (*pb.Empty, error) {
	clientID := domain.ClientID(h.contextProvider.MustAuthUID(ctx))
	app := h.clientApplication.BuildAsClient(clientID)

	if _, err := app.UpdateRole(
		ctx,
		domain.ClientID(req.ClientID),
		req.Role,
	); err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return &pb.Empty{}, nil
}

func (h *clientHandler) Delete(ctx context.Context, req *pb.ClientID) (*pb.Empty, error) {
	clientID := domain.ClientID(h.contextProvider.MustAuthUID(ctx))
	app := h.clientApplication.BuildAsClient(clientID)

	if err := app.Delete(
		ctx,
		domain.ClientID(req.Id),
	); err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return &pb.Empty{}, nil
}

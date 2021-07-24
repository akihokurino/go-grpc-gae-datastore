package api

import (
	"context"
	"time"

	"gae-go-sample/adapter"
	"gae-go-sample/domain"
	"gae-go-sample/handler"
	pb "gae-go-sample/proto/go/pb"
)

type entryHandler struct {
	errorConverter   adapter.ErrorConverter
	entryApplication adapter.EntryApplication
	contextProvider  handler.ContextProvider
}

func NewEntryHandler(
	errorConverter adapter.ErrorConverter,
	entryApplication adapter.EntryApplication,
	contextProvider handler.ContextProvider) pb.EntryService {
	return &entryHandler{
		errorConverter:   errorConverter,
		entryApplication: entryApplication,
		contextProvider:  contextProvider,
	}
}

func (h *entryHandler) GetAllByProject(ctx context.Context, req *pb.ProjectID) (*pb.EntryList, error) {
	clientID := domain.ClientID(h.contextProvider.MustAuthUID(ctx))
	app := h.entryApplication.BuildAsClient(clientID)

	entries, err := app.GetAllByProject(ctx, domain.ProjectID(req.Id))
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	response := make([]*pb.Entry, 0, len(entries))
	for _, entry := range entries {
		response = append(response, handler.ToEntryResponse(entry))
	}

	return &pb.EntryList{
		Items: response,
	}, nil
}

func (h *entryHandler) GetCountByProjects(ctx context.Context, req *pb.ProjectIDList) (*pb.EntryCountByProjectList, error) {
	clientID := domain.ClientID(h.contextProvider.MustAuthUID(ctx))
	app := h.entryApplication.BuildAsClient(clientID)

	projectIDs := make([]domain.ProjectID, 0, len(req.Ids))
	for _, id := range req.Ids {
		projectIDs = append(projectIDs, domain.ProjectID(id))
	}

	countMap, err := app.GetCountByProjects(ctx, projectIDs)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	responses := make([]*pb.EntryCountByProject, 0)
	for key, count := range countMap {
		responses = append(responses, &pb.EntryCountByProject{
			ProjectID: string(key),
			Count:     count,
		})
	}

	return &pb.EntryCountByProjectList{
		Items: responses,
	}, nil
}

func (h *entryHandler) GetCountByProject(ctx context.Context, req *pb.ProjectID) (*pb.Count, error) {
	clientID := domain.ClientID(h.contextProvider.MustAuthUID(ctx))
	app := h.entryApplication.BuildAsClient(clientID)

	count, err := app.GetCountByProject(ctx, domain.ProjectID(req.Id))
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return &pb.Count{Count: count}, nil
}

func (h *entryHandler) Create(ctx context.Context, req *pb.ProjectID) (*pb.Empty, error) {
	customerID := domain.CustomerID(h.contextProvider.MustAuthUID(ctx))
	app := h.entryApplication.BuildAsCustomer(customerID)

	now := time.Now()

	if _, err := app.Create(ctx, domain.ProjectID(req.Id), now); err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return &pb.Empty{}, nil
}

func (h *entryHandler) Delete(ctx context.Context, req *pb.ProjectID) (*pb.Empty, error) {
	customerID := domain.CustomerID(h.contextProvider.MustAuthUID(ctx))
	app := h.entryApplication.BuildAsCustomer(customerID)

	if err := app.Delete(ctx, domain.ProjectID(req.Id)); err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return &pb.Empty{}, nil
}

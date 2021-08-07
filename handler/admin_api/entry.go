package admin_api

import (
	"context"

	"gae-go-recruiting-server/adapter"
	"gae-go-recruiting-server/domain"
	"gae-go-recruiting-server/handler"
	pb "gae-go-recruiting-server/proto/go/pb"
)

type entryHandler struct {
	errorConverter   adapter.ErrorConverter
	entryApplication adapter.EntryApplication
	contextProvider  handler.ContextProvider
}

func NewEntryHandler(
	errorConverter adapter.ErrorConverter,
	entryApplication adapter.EntryApplication,
	contextProvider handler.ContextProvider) pb.AdminEntryService {
	return &entryHandler{
		errorConverter:   errorConverter,
		entryApplication: entryApplication,
		contextProvider:  contextProvider,
	}
}

func (h *entryHandler) GetAll(ctx context.Context, req *pb.Pager) (*pb.EntryList, error) {
	username := h.contextProvider.MustAuthAdminUID(ctx)
	app := h.entryApplication.BuildAsAdmin(username)

	entries, err := app.GetAllWithPager(ctx, req.Page, req.Offset)
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

func (h *entryHandler) GetAllByProject(ctx context.Context, req *pb.ProjectIDWithPager) (*pb.EntryList, error) {
	username := h.contextProvider.MustAuthAdminUID(ctx)
	app := h.entryApplication.BuildAsAdmin(username)

	entries, err := app.GetAllByProjectWithPager(ctx, domain.ProjectID(req.Id), req.Pager.Page, req.Pager.Offset)
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

func (h *entryHandler) Get(ctx context.Context, req *pb.EntryID) (*pb.Entry, error) {
	username := h.contextProvider.MustAuthAdminUID(ctx)
	app := h.entryApplication.BuildAsAdmin(username)

	entry, err := app.Get(ctx, domain.CustomerID(req.CustomerID), domain.ProjectID(req.ProjectID))
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return handler.ToEntryResponse(entry), nil
}

func (h *entryHandler) GetTotalCount(ctx context.Context, req *pb.Empty) (*pb.Count, error) {
	username := h.contextProvider.MustAuthAdminUID(ctx)
	app := h.entryApplication.BuildAsAdmin(username)

	count, err := app.GetTotalCount(ctx)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return &pb.Count{Count: count}, nil
}

func (h *entryHandler) GetCountByProject(ctx context.Context, req *pb.ProjectID) (*pb.Count, error) {
	username := h.contextProvider.MustAuthAdminUID(ctx)
	app := h.entryApplication.BuildAsAdmin(username)

	count, err := app.GetCountByProject(ctx, domain.ProjectID(req.Id))
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return &pb.Count{Count: count}, nil
}

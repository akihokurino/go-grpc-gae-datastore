package admin_api

import (
	"context"

	"gae-go-sample/adapter"
	"gae-go-sample/domain"
	"gae-go-sample/handler"
	pb "gae-go-sample/proto/go/pb"
)

type projectHandler struct {
	errorConverter     adapter.ErrorConverter
	projectApplication adapter.ProjectApplication
	contextProvider    handler.ContextProvider
}

func NewProjectHandler(
	errorConverter adapter.ErrorConverter,
	projectApplication adapter.ProjectApplication,
	contextProvider handler.ContextProvider) pb.AdminProjectService {
	return &projectHandler{
		errorConverter:     errorConverter,
		projectApplication: projectApplication,
		contextProvider:    contextProvider,
	}
}

func (h *projectHandler) GetAll(ctx context.Context, req *pb.Pager) (*pb.ProjectList, error) {
	username := h.contextProvider.MustAuthAdminUID(ctx)
	app := h.projectApplication.BuildAsAdmin(username)

	projects, err := app.GetAllWithPager(ctx, req.Page, req.Offset)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	responses := make([]*pb.Project, 0, len(projects))
	for _, project := range projects {
		responses = append(responses, handler.ToProjectResponse(project))
	}

	return &pb.ProjectList{
		Items: responses,
	}, nil
}

func (h *projectHandler) GetAllByIDs(ctx context.Context, req *pb.ProjectIDList) (*pb.ProjectList, error) {
	username := h.contextProvider.MustAuthAdminUID(ctx)
	app := h.projectApplication.BuildAsAdmin(username)

	ids := make([]domain.ProjectID, 0, len(req.Ids))
	for _, id := range req.Ids {
		ids = append(ids, domain.ProjectID(id))
	}

	projects, err := app.GetAllByIDs(ctx, ids)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	responses := make([]*pb.Project, 0, len(projects))
	for _, project := range projects {
		responses = append(responses, handler.ToProjectResponse(project))
	}

	return &pb.ProjectList{
		Items: responses,
	}, nil
}

func (h *projectHandler) GetAllByFilter(ctx context.Context, req *pb.FilterProjectRequest) (*pb.ProjectList, error) {
	username := h.contextProvider.MustAuthAdminUID(ctx)
	app := h.projectApplication.BuildAsAdmin(username)

	projects, err := app.GetAllByFilterWithPager(ctx, adapter.FilterProjectParams{
		Query: req.Q,
	}, req.Pager.Page, req.Pager.Offset)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	responses := make([]*pb.Project, 0, len(projects))
	for _, project := range projects {
		responses = append(responses, handler.ToProjectResponse(project))
	}

	return &pb.ProjectList{
		Items: responses,
	}, nil
}

func (h *projectHandler) GetAllByCompany(ctx context.Context, req *pb.CompanyIDWithPager) (*pb.ProjectList, error) {
	username := h.contextProvider.MustAuthAdminUID(ctx)
	app := h.projectApplication.BuildAsAdmin(username)

	projects, err := app.GetAllByCompanyWithPager(ctx, domain.CompanyID(req.Id), req.Pager.Page, req.Pager.Offset)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	responses := make([]*pb.Project, 0, len(projects))
	for _, project := range projects {
		responses = append(responses, handler.ToProjectResponse(project))
	}

	return &pb.ProjectList{
		Items: responses,
	}, nil
}

func (h *projectHandler) Get(ctx context.Context, req *pb.ProjectID) (*pb.Project, error) {
	username := h.contextProvider.MustAuthAdminUID(ctx)
	app := h.projectApplication.BuildAsAdmin(username)

	project, err := app.Get(ctx, domain.ProjectID(req.Id))
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return handler.ToProjectResponse(project), nil
}

func (h *projectHandler) GetTotalCount(ctx context.Context, req *pb.Empty) (*pb.Count, error) {
	username := h.contextProvider.MustAuthAdminUID(ctx)
	app := h.projectApplication.BuildAsAdmin(username)

	count, err := app.GetTotalCount(ctx)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return &pb.Count{Count: count}, nil
}

func (h *projectHandler) GetCountByFilter(ctx context.Context, req *pb.FilterProjectCountRequest) (*pb.Count, error) {
	username := h.contextProvider.MustAuthAdminUID(ctx)
	app := h.projectApplication.BuildAsAdmin(username)

	count, err := app.GetCountByFilter(ctx, adapter.FilterProjectParams{
		Query: req.Q,
	})
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return &pb.Count{Count: count}, nil
}

func (h *projectHandler) GetCountByCompany(ctx context.Context, req *pb.CompanyID) (*pb.Count, error) {
	username := h.contextProvider.MustAuthAdminUID(ctx)
	app := h.projectApplication.BuildAsAdmin(username)

	count, err := app.GetCountByCompany(ctx, domain.CompanyID(req.Id))
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return &pb.Count{Count: count}, nil
}

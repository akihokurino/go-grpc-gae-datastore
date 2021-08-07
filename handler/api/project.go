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

type projectHandler struct {
	errorConverter     adapter.ErrorConverter
	projectApplication adapter.ProjectApplication
	userApplication    adapter.UserApplication
	contextProvider    handler.ContextProvider
}

func NewProjectHandler(
	errorConverter adapter.ErrorConverter,
	projectApplication adapter.ProjectApplication,
	userApplication adapter.UserApplication,
	contextProvider handler.ContextProvider) pb.ProjectService {
	return &projectHandler{
		errorConverter:     errorConverter,
		projectApplication: projectApplication,
		userApplication:    userApplication,
		contextProvider:    contextProvider,
	}
}

func (h *projectHandler) GetAllByIDs(ctx context.Context, req *pb.ProjectIDList) (*pb.ProjectList, error) {
	userID := h.contextProvider.MustAuthUID(ctx)

	user, err := h.userApplication.Get(ctx, userID)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	ids := make([]domain.ProjectID, 0, len(req.Ids))
	for _, id := range req.Ids {
		ids = append(ids, domain.ProjectID(id))
	}

	var projects []*domain.Project

	switch user.Role {
	case pb.User_Role_Customer:
		app := h.projectApplication.BuildAsCustomer(user.CustomerID())

		projects, err = app.GetAllByIDs(ctx, ids)
	case pb.User_Role_Client:
		app := h.projectApplication.BuildAsClient(user.ClientID())

		projects, err = app.GetAllByIDs(ctx, ids)
	default:
		return nil, h.errorConverter(ctx, domain.ErrInvalidUserRole)
	}

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

func (h *projectHandler) GetAllByNewlyArrived(ctx context.Context, req *pb.Pager) (*pb.ProjectList, error) {
	customerID := domain.CustomerID(h.contextProvider.MustAuthUID(ctx))
	app := h.projectApplication.BuildAsCustomer(customerID)

	projects, err := app.GetAllByOpenWithPager(ctx, req.Page, req.Offset)
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

func (h *projectHandler) GetAllByRecommend(ctx context.Context, req *pb.Empty) (*pb.ProjectList, error) {
	customerID := domain.CustomerID(h.contextProvider.MustAuthUID(ctx))
	app := h.projectApplication.BuildAsCustomer(customerID)

	projects, err := app.GetAllByRecommend(ctx)
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

func (h *projectHandler) GetAllBySearch(ctx context.Context, req *pb.SearchProjectRequest) (*pb.ProjectList, error) {
	customerID := domain.CustomerID(h.contextProvider.MustAuthUID(ctx))
	app := h.projectApplication.BuildAsCustomer(customerID)

	projects, err := app.SearchWithPager(ctx, adapter.SearchProjectParams{
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

func (h *projectHandler) GetAllByClient(ctx context.Context, req *pb.Pager) (*pb.ProjectList, error) {
	clientID := domain.ClientID(h.contextProvider.MustAuthUID(ctx))
	app := h.projectApplication.BuildAsClient(clientID)

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

func (h *projectHandler) GetAllByClientExcludeAlreadyEntry(ctx context.Context, req *pb.CustomerID) (*pb.ProjectList, error) {
	clientID := domain.ClientID(h.contextProvider.MustAuthUID(ctx))
	app := h.projectApplication.BuildAsClient(clientID)

	customerID := domain.CustomerID(req.Id)

	projects, err := app.GetAllByOpenExcludeAlreadyEntry(ctx, customerID)
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

func (h *projectHandler) GetAllByEntry(ctx context.Context, req *pb.Pager) (*pb.ProjectList, error) {
	customerID := domain.CustomerID(h.contextProvider.MustAuthUID(ctx))
	app := h.projectApplication.BuildAsCustomer(customerID)

	projects, err := app.GetAllByEntryWithPager(ctx, req.Page, req.Offset)
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

func (h *projectHandler) GetAllByEntryNumOrder(ctx context.Context, req *pb.Pager) (*pb.ProjectList, error) {
	customerID := domain.CustomerID(h.contextProvider.MustAuthUID(ctx))
	app := h.projectApplication.BuildAsCustomer(customerID)

	projects, err := app.GetAllByEntryNumOrderWithPager(ctx, req.Page, req.Offset)
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
	userID := h.contextProvider.MustAuthUID(ctx)

	user, err := h.userApplication.Get(ctx, userID)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	var project *domain.Project

	switch user.Role {
	case pb.User_Role_Customer:
		app := h.projectApplication.BuildAsCustomer(user.CustomerID())

		project, err = app.Get(ctx, domain.ProjectID(req.Id))
	case pb.User_Role_Client:
		app := h.projectApplication.BuildAsClient(user.ClientID())

		project, err = app.Get(ctx, domain.ProjectID(req.Id))
	default:
		return nil, h.errorConverter(ctx, domain.ErrInvalidUserRole)
	}

	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return handler.ToProjectResponse(project), nil
}

func (h *projectHandler) GetCountByNewlyArrived(ctx context.Context, req *pb.Empty) (*pb.Count, error) {
	customerID := domain.CustomerID(h.contextProvider.MustAuthUID(ctx))
	app := h.projectApplication.BuildAsCustomer(customerID)

	count, err := app.GetCountByOpen(ctx)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return &pb.Count{Count: count}, nil
}

func (h *projectHandler) GetCountBySearch(ctx context.Context, req *pb.SearchProjectCountRequest) (*pb.Count, error) {
	customerID := domain.CustomerID(h.contextProvider.MustAuthUID(ctx))
	app := h.projectApplication.BuildAsCustomer(customerID)

	count, err := app.GetCountBySearch(ctx, adapter.SearchProjectParams{
		Query: req.Q,
	})
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return &pb.Count{Count: count}, nil
}

func (h *projectHandler) GetCountByClient(ctx context.Context, req *pb.Empty) (*pb.Count, error) {
	clientID := domain.ClientID(h.contextProvider.MustAuthUID(ctx))
	app := h.projectApplication.BuildAsClient(clientID)

	count, err := app.GetCount(ctx)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return &pb.Count{Count: count}, nil
}

func (h *projectHandler) GetCountByEntry(ctx context.Context, req *pb.Empty) (*pb.Count, error) {
	customerID := domain.CustomerID(h.contextProvider.MustAuthUID(ctx))
	app := h.projectApplication.BuildAsCustomer(customerID)

	count, err := app.GetCountByEntry(ctx)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return &pb.Count{Count: count}, nil
}

func (h *projectHandler) Create(ctx context.Context, req *pb.CreateProjectRequest) (*pb.Project, error) {
	clientID := domain.ClientID(h.contextProvider.MustAuthUID(ctx))
	app := h.projectApplication.BuildAsClient(clientID)

	now := time.Now()

	thumbnailURL, _ := url.Parse(req.ThumbnailURL)

	project, err := app.Create(
		ctx,
		adapter.ProjectParams{
			Name:         req.Name,
			Description:  req.Description,
			ThumbnailURL: thumbnailURL,
		},
		now)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return handler.ToProjectResponse(project), nil
}

func (h *projectHandler) Update(ctx context.Context, req *pb.UpdateProjectRequest) (*pb.Project, error) {
	clientID := domain.ClientID(h.contextProvider.MustAuthUID(ctx))
	app := h.projectApplication.BuildAsClient(clientID)

	thumbnailURL, _ := url.Parse(req.ThumbnailURL)

	project, err := app.Update(
		ctx,
		domain.ProjectID(req.Id),
		adapter.ProjectParams{
			Name:         req.Name,
			Description:  req.Description,
			ThumbnailURL: thumbnailURL,
		})
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return handler.ToProjectResponse(project), nil
}

func (h *projectHandler) Open(ctx context.Context, req *pb.ProjectID) (*pb.Empty, error) {
	clientID := domain.ClientID(h.contextProvider.MustAuthUID(ctx))
	app := h.projectApplication.BuildAsClient(clientID)

	now := time.Now()

	if err := app.Open(ctx, domain.ProjectID(req.Id), now); err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return &pb.Empty{}, nil
}

func (h *projectHandler) Draft(ctx context.Context, req *pb.ProjectID) (*pb.Empty, error) {
	clientID := domain.ClientID(h.contextProvider.MustAuthUID(ctx))
	app := h.projectApplication.BuildAsClient(clientID)

	if err := app.Draft(ctx, domain.ProjectID(req.Id)); err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return &pb.Empty{}, nil
}

func (h *projectHandler) Close(ctx context.Context, req *pb.ProjectID) (*pb.Empty, error) {
	clientID := domain.ClientID(h.contextProvider.MustAuthUID(ctx))
	app := h.projectApplication.BuildAsClient(clientID)

	if err := app.Close(ctx, domain.ProjectID(req.Id)); err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return &pb.Empty{}, nil
}

func (h *projectHandler) Delete(ctx context.Context, req *pb.ProjectID) (*pb.Empty, error) {
	clientID := domain.ClientID(h.contextProvider.MustAuthUID(ctx))
	app := h.projectApplication.BuildAsClient(clientID)

	if err := app.Delete(ctx, domain.ProjectID(req.Id)); err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return &pb.Empty{}, nil
}

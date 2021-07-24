package admin_api

import (
	"context"

	"gae-go-sample/adapter"
	"gae-go-sample/domain"
	"gae-go-sample/handler"
	pb "gae-go-sample/proto/go/pb"
)

type customerHandler struct {
	errorConverter      adapter.ErrorConverter
	customerApplication adapter.CustomerApplication
	contextProvider     handler.ContextProvider
}

func NewCustomerHandler(
	errorConverter adapter.ErrorConverter,
	customerApplication adapter.CustomerApplication,
	contextProvider handler.ContextProvider) pb.AdminCustomerService {
	return &customerHandler{
		errorConverter:      errorConverter,
		customerApplication: customerApplication,
		contextProvider:     contextProvider,
	}
}

func (h *customerHandler) GetAll(ctx context.Context, req *pb.Pager) (*pb.CustomerList, error) {
	username := h.contextProvider.MustAuthAdminUID(ctx)
	app := h.customerApplication.BuildAsAdmin(username)

	customers, err := app.GetAllWithPager(ctx, req.Page, req.Offset)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	responses := make([]*pb.Customer, 0, len(customers))
	for _, customer := range customers {
		responses = append(responses, handler.ToCustomerResponse(customer))
	}

	return &pb.CustomerList{
		Items: responses,
	}, nil
}

func (h *customerHandler) GetAllByIDs(ctx context.Context, req *pb.CustomerIDList) (*pb.CustomerList, error) {
	username := h.contextProvider.MustAuthAdminUID(ctx)
	app := h.customerApplication.BuildAsAdmin(username)

	ids := make([]domain.CustomerID, 0, len(req.Ids))
	for _, id := range req.Ids {
		ids = append(ids, domain.CustomerID(id))
	}

	customers, err := app.GetAllByIDs(ctx, ids)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	responses := make([]*pb.Customer, 0, len(customers))
	for _, customer := range customers {
		responses = append(responses, handler.ToCustomerResponse(customer))
	}

	return &pb.CustomerList{
		Items: responses,
	}, nil
}

func (h *customerHandler) GetAllByFilter(ctx context.Context, req *pb.FilterCustomerRequest) (*pb.CustomerList, error) {
	username := h.contextProvider.MustAuthAdminUID(ctx)
	app := h.customerApplication.BuildAsAdmin(username)

	customers, err := app.GetAllByFilterWithPager(ctx, adapter.FilterCustomerParams{
		Query:  req.Query,
		Status: req.Status,
	}, req.Pager.Page, req.Pager.Offset)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	responses := make([]*pb.Customer, 0, len(customers))
	for _, customer := range customers {
		responses = append(responses, handler.ToCustomerResponse(customer))
	}

	return &pb.CustomerList{
		Items: responses,
	}, nil
}

func (h *customerHandler) Get(ctx context.Context, req *pb.CustomerID) (*pb.Customer, error) {
	username := h.contextProvider.MustAuthAdminUID(ctx)
	app := h.customerApplication.BuildAsAdmin(username)

	customer, err := app.Get(ctx, domain.CustomerID(req.Id))
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return handler.ToCustomerResponse(customer), nil
}

func (h *customerHandler) GetTotalCount(ctx context.Context, req *pb.Empty) (*pb.Count, error) {
	username := h.contextProvider.MustAuthAdminUID(ctx)
	app := h.customerApplication.BuildAsAdmin(username)

	count, err := app.GetTotalCount(ctx)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return &pb.Count{Count: count}, nil
}

func (h *customerHandler) GetCountByFilter(ctx context.Context, req *pb.FilterCustomerCountRequest) (*pb.Count, error) {
	username := h.contextProvider.MustAuthAdminUID(ctx)
	app := h.customerApplication.BuildAsAdmin(username)

	count, err := app.GetCountByFilter(ctx, adapter.FilterCustomerParams{
		Query:  req.Query,
		Status: req.Status,
	})
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return &pb.Count{Count: count}, nil
}

func (h *customerHandler) PassInspection(ctx context.Context, req *pb.CustomerID) (*pb.Empty, error) {
	username := h.contextProvider.MustAuthAdminUID(ctx)
	app := h.customerApplication.BuildAsAdmin(username)

	if err := app.PassInspection(ctx, domain.CustomerID(req.Id)); err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return &pb.Empty{}, nil
}

func (h *customerHandler) DenyInspection(ctx context.Context, req *pb.CustomerID) (*pb.Empty, error) {
	username := h.contextProvider.MustAuthAdminUID(ctx)
	app := h.customerApplication.BuildAsAdmin(username)

	if err := app.DenyInspection(ctx, domain.CustomerID(req.Id)); err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return &pb.Empty{}, nil
}

func (h *customerHandler) Deny(ctx context.Context, req *pb.CustomerID) (*pb.Empty, error) {
	username := h.contextProvider.MustAuthAdminUID(ctx)
	app := h.customerApplication.BuildAsAdmin(username)

	if err := app.Deny(ctx, domain.CustomerID(req.Id)); err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return &pb.Empty{}, nil
}

func (h *customerHandler) ReInspect(ctx context.Context, req *pb.CustomerID) (*pb.Empty, error) {
	username := h.contextProvider.MustAuthAdminUID(ctx)
	app := h.customerApplication.BuildAsAdmin(username)

	if err := app.ReInspect(ctx, domain.CustomerID(req.Id)); err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return &pb.Empty{}, nil
}

func (h *customerHandler) Active(ctx context.Context, req *pb.CustomerID) (*pb.Empty, error) {
	username := h.contextProvider.MustAuthAdminUID(ctx)
	app := h.customerApplication.BuildAsAdmin(username)

	if err := app.Active(ctx, domain.CustomerID(req.Id)); err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return &pb.Empty{}, nil
}

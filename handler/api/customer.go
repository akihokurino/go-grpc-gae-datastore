package api

import (
	"context"
	"net/url"
	"time"

	"gae-go-sample/adapter"
	"gae-go-sample/domain"
	"gae-go-sample/handler"
	pb "gae-go-sample/proto/go/pb"
)

type customerHandler struct {
	errorConverter      adapter.ErrorConverter
	userApplication     adapter.UserApplication
	customerApplication adapter.CustomerApplication
	contextProvider     handler.ContextProvider
}

func NewCustomerHandler(
	errorConverter adapter.ErrorConverter,
	userApplication adapter.UserApplication,
	customerApplication adapter.CustomerApplication,
	contextProvider handler.ContextProvider) pb.CustomerService {
	return &customerHandler{
		errorConverter:      errorConverter,
		userApplication:     userApplication,
		customerApplication: customerApplication,
		contextProvider:     contextProvider,
	}
}

func (h *customerHandler) GetAll(ctx context.Context, req *pb.Empty) (*pb.CustomerList, error) {
	userID := h.contextProvider.MustAuthUID(ctx)

	user, err := h.userApplication.Get(ctx, userID)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	switch user.Role {
	case pb.User_Role_Client:
		app := h.customerApplication.BuildAsClient(user.ClientID())

		customers, err := app.GetAll(ctx)
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
	default:
		return nil, h.errorConverter(ctx, domain.ErrInvalidUserRole)
	}
}

func (h *customerHandler) GetAllByIDs(ctx context.Context, req *pb.CustomerIDList) (*pb.CustomerList, error) {
	userID := h.contextProvider.MustAuthUID(ctx)

	user, err := h.userApplication.Get(ctx, userID)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	switch user.Role {
	case pb.User_Role_Client:
		app := h.customerApplication.BuildAsClient(user.ClientID())

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
	default:
		return nil, h.errorConverter(ctx, domain.ErrInvalidUserRole)
	}
}

func (h *customerHandler) GetAllBySearch(ctx context.Context, req *pb.SearchCustomerRequest) (*pb.CustomerList, error) {
	clientID := domain.ClientID(h.contextProvider.MustAuthUID(ctx))
	userID := h.contextProvider.MustAuthUID(ctx)

	user, err := h.userApplication.Get(ctx, userID)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	switch user.Role {
	case pb.User_Role_Client:
		app := h.customerApplication.BuildAsClient(clientID)

		customers, err := app.SearchWithPager(ctx, adapter.SearchCustomerParams{
			Query:   req.Query,
			OrderBy: req.OrderBy,
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
	default:
		return nil, h.errorConverter(ctx, domain.ErrInvalidUserRole)
	}
}

func (h *customerHandler) Get(ctx context.Context, req *pb.CustomerID) (*pb.Customer, error) {
	clientID := domain.ClientID(h.contextProvider.MustAuthUID(ctx))
	app := h.customerApplication.BuildAsClient(clientID)

	customer, err := app.Get(ctx, domain.CustomerID(req.Id))
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return handler.ToCustomerResponse(customer), nil
}

func (h *customerHandler) GetCountBySearch(ctx context.Context, req *pb.SearchCustomerCountRequest) (*pb.Count, error) {
	userID := h.contextProvider.MustAuthUID(ctx)

	user, err := h.userApplication.Get(ctx, userID)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	switch user.Role {
	case pb.User_Role_Client:
		app := h.customerApplication.BuildAsClient(user.ClientID())

		count, err := app.GetCountBySearch(ctx, adapter.SearchCustomerParams{
			Query: req.Query,
		})
		if err != nil {
			return nil, h.errorConverter(ctx, err)
		}

		return &pb.Count{Count: count}, nil
	default:
		return nil, h.errorConverter(ctx, domain.ErrInvalidUserRole)
	}
}

func (h *customerHandler) Create(ctx context.Context, req *pb.CreateCustomerRequest) (*pb.Customer, error) {
	customerID := domain.CustomerID(h.contextProvider.MustAuthUID(ctx))
	app := h.customerApplication.BuildAsCustomer(customerID)

	now := time.Now()
	birthdate, _ := time.Parse("2006-01-02", req.Birthdate)

	customer, err := app.Create(
		ctx,
		req.Name,
		req.NameKana,
		req.Gender,
		req.PhoneNumber,
		birthdate,
		now,
	)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return handler.ToCustomerResponse(customer), nil
}

func (h *customerHandler) Update(ctx context.Context, req *pb.UpdateCustomerRequest) (*pb.Customer, error) {
	customerID := domain.CustomerID(h.contextProvider.MustAuthUID(ctx))
	app := h.customerApplication.BuildAsCustomer(customerID)

	iconURL, _ := url.Parse(req.IconURL)
	resumeURL, _ := url.Parse(req.ResumeURL)

	birthdate, _ := time.Parse("2006-01-02", req.Birthdate)

	customer, err := app.Update(
		ctx,
		adapter.CustomerParams{
			Name:        req.Name,
			NameKana:    req.NameKana,
			IconURL:     iconURL,
			Birthdate:   birthdate,
			Gender:      req.Gender,
			PhoneNumber: req.PhoneNumber,
			Pr:          req.Pr,
			Address:     req.Address,
			ResumeURL:   resumeURL,
		},
	)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return handler.ToCustomerResponse(customer), nil
}

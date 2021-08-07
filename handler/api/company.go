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

type companyHandler struct {
	errorConverter     adapter.ErrorConverter
	companyApplication adapter.CompanyApplication
	userApplication    adapter.UserApplication
	contextProvider    handler.ContextProvider
}

func NewCompanyHandler(
	errorConverter adapter.ErrorConverter,
	companyApplication adapter.CompanyApplication,
	userApplication adapter.UserApplication,
	contextProvider handler.ContextProvider) pb.CompanyService {
	return &companyHandler{
		errorConverter:     errorConverter,
		companyApplication: companyApplication,
		userApplication:    userApplication,
		contextProvider:    contextProvider,
	}
}

func (h *companyHandler) GetAllByIDs(ctx context.Context, req *pb.CompanyIDList) (*pb.CompanyList, error) {
	userID := h.contextProvider.MustAuthUID(ctx)

	ids := make([]domain.CompanyID, 0, len(req.Ids))
	for _, id := range req.Ids {
		ids = append(ids, domain.CompanyID(id))
	}

	user, err := h.userApplication.Get(ctx, userID)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	var companies []*domain.Company

	switch user.Role {
	case pb.User_Role_Customer:
		app := h.companyApplication.BuildAsCustomer(domain.CustomerID(userID))

		companies, err = app.GetAllByIDsWithMaskIfNotActive(ctx, ids)
	case pb.User_Role_Client:
		app := h.companyApplication.BuildAsClient(domain.ClientID(userID))

		companies, err = app.GetAllByIDs(ctx, ids)
	default:
		return nil, h.errorConverter(ctx, domain.ErrInvalidUserRole)
	}

	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	responses := make([]*pb.Company, 0, len(companies))

	for _, company := range companies {
		responses = append(responses, handler.ToCompanyResponse(company))
	}

	return &pb.CompanyList{Items: responses}, nil
}

func (h *companyHandler) GetAllByIDsWithMask(ctx context.Context, req *pb.CompanyIDWithProjectIDList) (*pb.CompanyWithProjectIDList, error) {
	userID := h.contextProvider.MustAuthUID(ctx)

	ids := make([]adapter.CompanyIDWithProjectID, 0, len(req.Items))
	for _, item := range req.Items {
		ids = append(ids, adapter.CompanyIDWithProjectID{
			CompanyID: domain.CompanyID(item.CompanyID),
			ProjectID: domain.ProjectID(item.ProjectID),
		})
	}

	user, err := h.userApplication.Get(ctx, userID)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	var companyWithProjectIDs []*adapter.CompanyWithProjectID

	switch user.Role {
	case pb.User_Role_Customer:
		app := h.companyApplication.BuildAsCustomer(domain.CustomerID(userID))

		companyWithProjectIDs, err = app.GetAllByIDsWithMaskIfNeedOrNotActive(ctx, ids)
	case pb.User_Role_Client:
		app := h.companyApplication.BuildAsClient(domain.ClientID(userID))

		companyWithProjectIDs, err = app.GetAllByIDsWithMaskIfNeed(ctx, ids)
	default:
		return nil, h.errorConverter(ctx, domain.ErrInvalidUserRole)
	}

	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	responses := make([]*pb.CompanyWithProjectID, 0, len(companyWithProjectIDs))

	for _, id := range companyWithProjectIDs {
		responses = append(responses, &pb.CompanyWithProjectID{
			Company:   handler.ToCompanyResponse(id.Company),
			ProjectID: string(id.ProjectID),
		})
	}

	return &pb.CompanyWithProjectIDList{Items: responses}, nil
}

func (h *companyHandler) Get(ctx context.Context, req *pb.CompanyID) (*pb.Company, error) {
	userID := h.contextProvider.MustAuthUID(ctx)

	user, err := h.userApplication.Get(ctx, userID)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	var company *domain.Company

	switch user.Role {
	case pb.User_Role_Customer:
		app := h.companyApplication.BuildAsCustomer(domain.CustomerID(userID))

		company, err = app.GetWithMaskIfNotActive(ctx, domain.CompanyID(req.Id))
	case pb.User_Role_Client:
		app := h.companyApplication.BuildAsClient(domain.ClientID(userID))

		company, err = app.Get(ctx, domain.CompanyID(req.Id))
	default:
		return nil, h.errorConverter(ctx, domain.ErrInvalidUserRole)
	}

	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return handler.ToCompanyResponse(company), nil
}

func (h *companyHandler) GetWithMask(ctx context.Context, req *pb.CompanyIDWithProjectID) (*pb.Company, error) {
	userID := h.contextProvider.MustAuthUID(ctx)

	user, err := h.userApplication.Get(ctx, userID)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	var company *domain.Company

	switch user.Role {
	case pb.User_Role_Customer:
		app := h.companyApplication.BuildAsCustomer(domain.CustomerID(userID))

		company, err = app.GetWithMaskIfNeedOrNotActive(
			ctx,
			adapter.CompanyIDWithProjectID{
				CompanyID: domain.CompanyID(req.CompanyID),
				ProjectID: domain.ProjectID(req.ProjectID),
			})
	case pb.User_Role_Client:
		app := h.companyApplication.BuildAsClient(domain.ClientID(userID))

		company, err = app.GetWithMaskIfNeed(ctx, adapter.CompanyIDWithProjectID{
			CompanyID: domain.CompanyID(req.CompanyID),
			ProjectID: domain.ProjectID(req.ProjectID),
		})
	default:
		return nil, h.errorConverter(ctx, domain.ErrInvalidUserRole)
	}

	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return handler.ToCompanyResponse(company), nil
}

func (h *companyHandler) Update(ctx context.Context, req *pb.UpdateCompanyRequest) (*pb.Company, error) {
	clientID := domain.ClientID(h.contextProvider.MustAuthUID(ctx))
	app := h.companyApplication.BuildAsClient(clientID)

	logoURL, _ := url.Parse(req.LogoURL)

	webURL, _ := url.Parse(req.WebURL)

	establishedAt, _ := time.Parse("2006-01-02", req.EstablishedAt)

	company, err := app.Update(
		ctx,
		adapter.CompanyParams{
			Name:                       req.Name,
			LogoURL:                    logoURL,
			WebURL:                     webURL,
			EstablishedAt:              establishedAt,
			PostalCode:                 req.PostalCode,
			RepresentativeName:         req.RepresentativeName,
			CapitalStock:               req.CapitalStock,
			Introduction:               req.Introduction,
			AccordingCompanyName:       req.AccordingCompanyName,
			AccordingCompanyPostalCode: req.AccordingCompanyPostalCode,
			AccordingCompanyAddress:    req.AccordingCompanyAddress,
		},
	)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return handler.ToCompanyResponse(company), nil
}

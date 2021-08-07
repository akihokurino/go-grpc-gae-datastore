package admin_api

import (
	"context"

	"gae-go-recruiting-server/adapter"
	"gae-go-recruiting-server/domain"
	"gae-go-recruiting-server/handler"
	pb "gae-go-recruiting-server/proto/go/pb"
)

type companyHandler struct {
	errorConverter     adapter.ErrorConverter
	companyApplication adapter.CompanyApplication
	contextProvider    handler.ContextProvider
}

func NewCompanyHandler(
	errorConverter adapter.ErrorConverter,
	companyApplication adapter.CompanyApplication,
	contextProvider handler.ContextProvider) pb.AdminCompanyService {
	return &companyHandler{
		errorConverter:     errorConverter,
		companyApplication: companyApplication,
		contextProvider:    contextProvider,
	}
}

func (h *companyHandler) GetAll(ctx context.Context, req *pb.Pager) (*pb.CompanyList, error) {
	username := h.contextProvider.MustAuthAdminUID(ctx)
	app := h.companyApplication.BuildAsAdmin(username)

	companies, err := app.GetAllWithPager(ctx, req.Page, req.Offset)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	responses := make([]*pb.Company, 0, len(companies))
	for _, company := range companies {
		responses = append(responses, handler.ToCompanyResponse(company))
	}

	return &pb.CompanyList{
		Items: responses,
	}, nil
}

func (h *companyHandler) GetAllByIDs(ctx context.Context, req *pb.CompanyIDList) (*pb.CompanyList, error) {
	username := h.contextProvider.MustAuthAdminUID(ctx)
	app := h.companyApplication.BuildAsAdmin(username)

	ids := make([]domain.CompanyID, 0, len(req.Ids))
	for _, id := range req.Ids {
		ids = append(ids, domain.CompanyID(id))
	}

	companies, err := app.GetAllByIDs(ctx, ids)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	responses := make([]*pb.Company, 0, len(companies))
	for _, company := range companies {
		responses = append(responses, handler.ToCompanyResponse(company))
	}

	return &pb.CompanyList{
		Items: responses,
	}, nil
}

func (h *companyHandler) Get(ctx context.Context, req *pb.CompanyID) (*pb.Company, error) {
	username := h.contextProvider.MustAuthAdminUID(ctx)
	app := h.companyApplication.BuildAsAdmin(username)

	company, err := app.Get(ctx, domain.CompanyID(req.Id))
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return handler.ToCompanyResponse(company), nil
}

func (h *companyHandler) GetTotalCount(ctx context.Context, req *pb.Empty) (*pb.Count, error) {
	username := h.contextProvider.MustAuthAdminUID(ctx)
	app := h.companyApplication.BuildAsAdmin(username)

	count, err := app.GetTotalCount(ctx)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return &pb.Count{Count: count}, nil
}

func (h *companyHandler) Active(ctx context.Context, req *pb.CompanyID) (*pb.Empty, error) {
	username := h.contextProvider.MustAuthAdminUID(ctx)
	app := h.companyApplication.BuildAsAdmin(username)

	if err := app.Active(ctx, domain.CompanyID(req.Id)); err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return &pb.Empty{}, nil
}

func (h *companyHandler) Ban(ctx context.Context, req *pb.CompanyID) (*pb.Empty, error) {
	username := h.contextProvider.MustAuthAdminUID(ctx)
	app := h.companyApplication.BuildAsAdmin(username)

	if err := app.Ban(ctx, domain.CompanyID(req.Id)); err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return &pb.Empty{}, nil
}

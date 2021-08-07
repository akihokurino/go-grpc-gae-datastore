package admin_api

import (
	"context"

	"gae-go-recruiting-server/adapter"
	"gae-go-recruiting-server/domain"
	"gae-go-recruiting-server/handler"
	pb "gae-go-recruiting-server/proto/go/pb"
)

type contractHandler struct {
	errorConverter      adapter.ErrorConverter
	contractApplication adapter.ContractApplication
	contextProvider     handler.ContextProvider
}

func NewContractHandler(
	errorConverter adapter.ErrorConverter,
	contractApplication adapter.ContractApplication,
	contextProvider handler.ContextProvider) pb.AdminContractService {
	return &contractHandler{
		errorConverter:      errorConverter,
		contractApplication: contractApplication,
		contextProvider:     contextProvider,
	}
}

func (h *contractHandler) GetAll(ctx context.Context, req *pb.Pager) (*pb.ContractList, error) {
	username := h.contextProvider.MustAuthAdminUID(ctx)
	app := h.contractApplication.BuildAsAdmin(username)

	contracts, err := app.GetAllWithPager(ctx, req.Page, req.Offset)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	responses := make([]*pb.Contract, 0, len(contracts))
	for _, contract := range contracts {
		responses = append(responses, handler.ToContractResponse(contract))
	}

	return &pb.ContractList{
		Items: responses,
	}, nil
}

func (h *contractHandler) Get(ctx context.Context, req *pb.ContractID) (*pb.Contract, error) {
	username := h.contextProvider.MustAuthAdminUID(ctx)
	app := h.contractApplication.BuildAsAdmin(username)

	contract, err := app.Get(ctx, adapter.ContractIDParams{
		CompanyID:  domain.CompanyID(req.CompanyID),
		ProjectID:  domain.ProjectID(req.ProjectID),
		CustomerID: domain.CustomerID(req.CustomerID),
	})
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return handler.ToContractResponse(contract), nil
}

func (h *contractHandler) GetTotalCount(ctx context.Context, req *pb.Empty) (*pb.Count, error) {
	username := h.contextProvider.MustAuthAdminUID(ctx)
	app := h.contractApplication.BuildAsAdmin(username)

	count, err := app.GetTotalCount(ctx)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return &pb.Count{Count: count}, nil
}

func (h *contractHandler) Accept(ctx context.Context, req *pb.ContractID) (*pb.Empty, error) {
	username := h.contextProvider.MustAuthAdminUID(ctx)
	app := h.contractApplication.BuildAsAdmin(username)

	if err := app.Accept(ctx, adapter.ContractIDParams{
		CompanyID:  domain.CompanyID(req.CompanyID),
		ProjectID:  domain.ProjectID(req.ProjectID),
		CustomerID: domain.CustomerID(req.CustomerID),
	}); err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return &pb.Empty{}, nil
}

func (h *contractHandler) Cancel(ctx context.Context, req *pb.ContractID) (*pb.Empty, error) {
	username := h.contextProvider.MustAuthAdminUID(ctx)
	app := h.contractApplication.BuildAsAdmin(username)

	if err := app.Cancel(ctx, adapter.ContractIDParams{
		CompanyID:  domain.CompanyID(req.CompanyID),
		ProjectID:  domain.ProjectID(req.ProjectID),
		CustomerID: domain.CustomerID(req.CustomerID),
	}); err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return &pb.Empty{}, nil
}

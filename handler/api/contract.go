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

type contractHandler struct {
	errorConverter      adapter.ErrorConverter
	contractApplication adapter.ContractApplication
	contextProvider     handler.ContextProvider
}

func NewContractHandler(
	errorConverter adapter.ErrorConverter,
	contractApplication adapter.ContractApplication,
	contextProvider handler.ContextProvider) pb.ContractService {
	return &contractHandler{
		errorConverter:      errorConverter,
		contractApplication: contractApplication,
		contextProvider:     contextProvider,
	}
}

func (h *contractHandler) GetAllNewestByIDs(ctx context.Context, req *pb.ContractIDWithoutCompanyIDList) (*pb.ContractList, error) {
	clientID := domain.ClientID(h.contextProvider.MustAuthUID(ctx))
	app := h.contractApplication.BuildAsClient(clientID)

	idParams := make([]adapter.ContractIDWithoutCompanyIDParams, 0, len(req.Ids))
	for _, id := range req.Ids {
		idParams = append(idParams, adapter.ContractIDWithoutCompanyIDParams{
			ProjectID:  domain.ProjectID(id.ProjectID),
			CustomerID: domain.CustomerID(id.CustomerID),
		})
	}

	contracts, err := app.GetAllNewestByIDs(ctx, idParams)
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

func (h *contractHandler) GetAllByClient(ctx context.Context, req *pb.Pager) (*pb.ContractList, error) {
	clientID := domain.ClientID(h.contextProvider.MustAuthUID(ctx))
	app := h.contractApplication.BuildAsClient(clientID)

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
	clientID := domain.ClientID(h.contextProvider.MustAuthUID(ctx))
	app := h.contractApplication.BuildAsClient(clientID)

	projectID := domain.ProjectID(req.ProjectID)
	customerID := domain.CustomerID(req.CustomerID)

	contract, err := app.Get(ctx, adapter.ContractIDWithoutCompanyIDParams{
		ProjectID:  projectID,
		CustomerID: customerID,
	})
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return handler.ToContractResponse(contract), nil
}

func (h *contractHandler) GetCountByClient(ctx context.Context, req *pb.Empty) (*pb.Count, error) {
	clientID := domain.ClientID(h.contextProvider.MustAuthUID(ctx))
	app := h.contractApplication.BuildAsClient(clientID)

	count, err := app.GetCount(ctx)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return &pb.Count{Count: count}, nil
}

func (h *contractHandler) Create(ctx context.Context, req *pb.CreateContractRequest) (*pb.Contract, error) {
	clientID := domain.ClientID(h.contextProvider.MustAuthUID(ctx))
	app := h.contractApplication.BuildAsClient(clientID)

	projectID := domain.ProjectID(req.ProjectID)
	customerID := domain.CustomerID(req.CustomerID)

	fileURL, err := url.Parse(req.FileURL)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	now := time.Now()

	contract, err := app.Create(ctx, adapter.ContractIDWithoutCompanyIDParams{
		ProjectID:  projectID,
		CustomerID: customerID,
	}, fileURL, now)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return handler.ToContractResponse(contract), nil
}

func (h *contractHandler) Update(ctx context.Context, req *pb.UpdateContractRequest) (*pb.Contract, error) {
	clientID := domain.ClientID(h.contextProvider.MustAuthUID(ctx))
	app := h.contractApplication.BuildAsClient(clientID)

	projectID := domain.ProjectID(req.ProjectID)
	customerID := domain.CustomerID(req.CustomerID)

	fileURL, err := url.Parse(req.FileURL)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	now := time.Now()

	contract, err := app.Update(ctx, adapter.ContractIDWithoutCompanyIDParams{
		ProjectID:  projectID,
		CustomerID: customerID,
	}, fileURL, now)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return handler.ToContractResponse(contract), nil
}

func (h *contractHandler) Delete(ctx context.Context, req *pb.DeleteContractRequest) (*pb.Empty, error) {
	clientID := domain.ClientID(h.contextProvider.MustAuthUID(ctx))
	app := h.contractApplication.BuildAsClient(clientID)

	projectID := domain.ProjectID(req.ProjectID)
	customerID := domain.CustomerID(req.CustomerID)

	if err := app.Delete(ctx, adapter.ContractIDWithoutCompanyIDParams{
		ProjectID:  projectID,
		CustomerID: customerID,
	}); err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return &pb.Empty{}, nil
}

package admin_api

import (
	"context"
	"time"

	"gae-go-sample/adapter"
	"gae-go-sample/domain"
	"gae-go-sample/handler"
	pb "gae-go-sample/proto/go/pb"
)

type applyClientHandler struct {
	errorConverter         adapter.ErrorConverter
	applyClientApplication adapter.ApplyClientApplication
	contextProvider        handler.ContextProvider
}

func NewApplyClientHandler(
	errorConverter adapter.ErrorConverter,
	applyClientApplication adapter.ApplyClientApplication,
	contextProvider handler.ContextProvider) pb.AdminApplyClientService {
	return &applyClientHandler{
		errorConverter:         errorConverter,
		applyClientApplication: applyClientApplication,
		contextProvider:        contextProvider,
	}
}

func (h *applyClientHandler) GetAll(ctx context.Context, req *pb.Pager) (*pb.ApplyClientList, error) {
	username := h.contextProvider.MustAuthAdminUID(ctx)
	app := h.applyClientApplication.BuildAsAdmin(username)

	applies, err := app.GetAllWithPager(ctx, req.Page, req.Offset)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	responses := make([]*pb.ApplyClient, 0, len(applies))
	for _, applyClient := range applies {
		responses = append(responses, handler.ToApplyClientResponse(applyClient))
	}

	return &pb.ApplyClientList{
		Items: responses,
	}, nil
}

func (h *applyClientHandler) GetAllByFilter(ctx context.Context, req *pb.FilterApplyClientRequest) (*pb.ApplyClientList, error) {
	username := h.contextProvider.MustAuthAdminUID(ctx)
	app := h.applyClientApplication.BuildAsAdmin(username)

	applies, err := app.GetAllByFilterWithPager(ctx, req.Status, req.Pager.Page, req.Pager.Offset)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	responses := make([]*pb.ApplyClient, 0, len(applies))
	for _, applyClient := range applies {
		responses = append(responses, handler.ToApplyClientResponse(applyClient))
	}

	return &pb.ApplyClientList{
		Items: responses,
	}, nil
}

func (h *applyClientHandler) Get(ctx context.Context, req *pb.ApplyClientID) (*pb.ApplyClient, error) {
	username := h.contextProvider.MustAuthAdminUID(ctx)
	app := h.applyClientApplication.BuildAsAdmin(username)

	apply, err := app.Get(ctx, domain.ApplyClientID(req.Email))
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return handler.ToApplyClientResponse(apply), nil
}

func (h *applyClientHandler) GetTotalCount(ctx context.Context, req *pb.Empty) (*pb.Count, error) {
	username := h.contextProvider.MustAuthAdminUID(ctx)
	app := h.applyClientApplication.BuildAsAdmin(username)

	count, err := app.GetTotalCount(ctx)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return &pb.Count{Count: count}, nil
}

func (h *applyClientHandler) GetCountByFilter(ctx context.Context, req *pb.FilterApplyClientCountRequest) (*pb.Count, error) {
	username := h.contextProvider.MustAuthAdminUID(ctx)
	app := h.applyClientApplication.BuildAsAdmin(username)

	count, err := app.GetCountByFilter(ctx, req.Status)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return &pb.Count{Count: count}, nil
}

func (h *applyClientHandler) Accept(ctx context.Context, req *pb.AcceptApplyClientRequest) (*pb.Empty, error) {
	now := time.Now()

	username := h.contextProvider.MustAuthAdminUID(ctx)
	app := h.applyClientApplication.BuildAsAdmin(username)

	if err := app.Accept(
		ctx,
		domain.ApplyClientID(req.Email),
		req.Password,
		now); err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return &pb.Empty{}, nil
}

func (h *applyClientHandler) Deny(ctx context.Context, req *pb.ApplyClientID) (*pb.Empty, error) {
	username := h.contextProvider.MustAuthAdminUID(ctx)
	app := h.applyClientApplication.BuildAsAdmin(username)

	if _, err := app.Deny(ctx, domain.ApplyClientID(req.Email)); err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return &pb.Empty{}, nil
}

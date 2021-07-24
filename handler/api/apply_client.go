package api

import (
	"context"
	"net/url"
	"time"

	"gae-go-sample/adapter"
	"gae-go-sample/domain"
	pb "gae-go-sample/proto/go/pb"
)

type applyClientHandler struct {
	errorConverter         adapter.ErrorConverter
	applyClientApplication adapter.ApplyClientApplication
}

func NewApplyClientHandler(
	errorConverter adapter.ErrorConverter,
	applyClientApplication adapter.ApplyClientApplication) pb.ApplyClientService {
	return &applyClientHandler{
		errorConverter:         errorConverter,
		applyClientApplication: applyClientApplication,
	}
}

func (h *applyClientHandler) Register(ctx context.Context, req *pb.RegisterClientRequest) (*pb.Empty, error) {
	now := time.Now()

	app := h.applyClientApplication.BuildAsPublic()

	webURL, err := url.Parse(req.WebURL)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	if _, err := app.Create(
		ctx,
		adapter.ApplyClientParams{
			Email:           domain.ApplyClientID(req.Email),
			PhoneNumber:     req.PhoneNumber,
			CompanyName:     req.CompanyName,
			WebURL:          webURL,
			AccountName:     req.AccountName,
			AccountNameKana: req.AccountNameKana,
			Position:        req.Position,
		},
		now); err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	return &pb.Empty{}, nil
}

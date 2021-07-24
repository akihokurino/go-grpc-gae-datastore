package api

import (
	"context"

	"gae-go-sample/adapter"
	"gae-go-sample/handler"
	pb "gae-go-sample/proto/go/pb"
)

type prefectureHandler struct {
	errorConverter        adapter.ErrorConverter
	prefectureApplication adapter.PrefectureApplication
}

func NewPrefectureHandler(
	errorConverter adapter.ErrorConverter,
	prefectureApplication adapter.PrefectureApplication) pb.PrefectureService {
	return &prefectureHandler{
		errorConverter:        errorConverter,
		prefectureApplication: prefectureApplication,
	}
}

func (h *prefectureHandler) GetAll(ctx context.Context, req *pb.Empty) (*pb.PrefectureList, error) {

	prefectures, err := h.prefectureApplication.GetAll(ctx)
	if err != nil {
		return nil, h.errorConverter(ctx, err)
	}

	responses := make([]*pb.Prefecture, 0, len(prefectures))
	for _, pref := range prefectures {
		responses = append(responses, handler.ToPrefectureResponse(pref))
	}

	return &pb.PrefectureList{
		Items: responses,
	}, nil
}

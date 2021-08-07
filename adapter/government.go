package adapter

import (
	"context"

	"gae-go-recruiting-server/domain"
)

type GovernmentClient interface {
	Get(ctx context.Context, pathStr string, params map[string]string) ([]byte, error)
}

type PrefectureRepository interface {
	GetAll(ctx context.Context) ([]*domain.Prefecture, error)
}

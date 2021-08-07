package adapter

import (
	"context"
	"net/url"

	"gae-go-recruiting-server/domain"
)

type GCSSignature func(ctx context.Context, gsURL *url.URL) (*url.URL, error)

type ContractFileRepository interface {
	Get(ctx context.Context, contract *domain.Contract) (*domain.File, error)
}

type AdminUserFileRepository interface {
	Save(ctx context.Context, username string, data []byte) (*url.URL, error)
}

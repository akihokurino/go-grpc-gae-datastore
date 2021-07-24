package prefecture

import (
	"context"
	"encoding/json"

	"gae-go-sample/adapter"
	"gae-go-sample/domain"

	"github.com/pkg/errors"
)

type repository struct {
	client adapter.GovernmentClient
}

func NewRepository(client adapter.GovernmentClient) adapter.PrefectureRepository {
	return &repository{
		client: client,
	}
}

func (r *repository) GetAll(ctx context.Context) ([]*domain.Prefecture, error) {
	bytes, err := r.client.Get(ctx, "/api/v1/prefectures", map[string]string{})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	result := root{}
	if err := json.Unmarshal(bytes, &result); err != nil {
		return nil, errors.WithStack(err)
	}

	items := make([]*domain.Prefecture, 0, len(result.Items))
	for _, item := range result.Items {
		items = append(items, item.toDomain())
	}

	return items, nil
}

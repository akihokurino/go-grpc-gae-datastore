package no_message_support_table

import (
	"context"

	"go.mercari.io/datastore/boom"

	"gae-go-sample/adapter"
	"gae-go-sample/domain"

	"github.com/pkg/errors"
)

func NewRepository(client adapter.DSClient) adapter.NoMessageSupportRepository {
	return &repository{
		client: client,
	}
}

type repository struct {
	client adapter.DSClient
}

func (r *repository) GetAllByOpened(ctx context.Context) ([]*domain.NoMessageSupport, error) {
	var entities []*entity

	if err := r.client.GetByFilter(
		ctx,
		kind,
		&entities,
		map[string]interface{}{
			"Closed =": false,
		},
		nil,
		""); err != nil {
		return nil, errors.WithStack(err)
	}

	items := make([]*domain.NoMessageSupport, 0, len(entities))
	for _, e := range entities {
		items = append(items, e.toDomain())
	}

	return items, nil
}

func (r *repository) Get(ctx context.Context, id domain.NoMessageSupportID) (*domain.NoMessageSupport, error) {
	entity := onlyID(id)

	if err := r.client.Get(ctx, entity); err != nil {
		return nil, errors.WithStack(err)
	}

	return entity.toDomain(), nil
}

func (r *repository) Put(tx *boom.Transaction, support *domain.NoMessageSupport) error {
	if err := r.client.Put(tx, toEntity(support)); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

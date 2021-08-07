package no_entry_support_table

import (
	"context"

	"go.mercari.io/datastore/boom"

	"gae-go-recruiting-server/adapter"
	"gae-go-recruiting-server/domain"

	"github.com/pkg/errors"
)

func NewRepository(client adapter.DSClient) adapter.NoEntrySupportRepository {
	return &repository{
		client: client,
	}
}

type repository struct {
	client adapter.DSClient
}

func (r *repository) GetAllByOpened(ctx context.Context) ([]*domain.NoEntrySupport, error) {
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

	items := make([]*domain.NoEntrySupport, 0, len(entities))
	for _, e := range entities {
		items = append(items, e.toDomain())
	}

	return items, nil
}

func (r *repository) Get(ctx context.Context, projectID domain.ProjectID) (*domain.NoEntrySupport, error) {
	entity := onlyID(projectID)

	if err := r.client.Get(ctx, entity); err != nil {
		return nil, errors.WithStack(err)
	}

	return entity.toDomain(), nil
}

func (r *repository) Put(tx *boom.Transaction, support *domain.NoEntrySupport) error {
	if err := r.client.Put(tx, toEntity(support)); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

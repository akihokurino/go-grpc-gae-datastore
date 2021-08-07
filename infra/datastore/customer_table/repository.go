package customer_table

import (
	"context"

	"go.mercari.io/datastore/boom"

	"gae-go-recruiting-server/adapter"
	"gae-go-recruiting-server/domain"

	"github.com/pkg/errors"
)

func NewRepository(client adapter.DSClient) adapter.CustomerRepository {
	return &repository{
		client: client,
	}
}

type repository struct {
	client adapter.DSClient
}

func (r *repository) GetAll(ctx context.Context) ([]*domain.Customer, error) {
	var entities []*entity

	if err := r.client.GetAll(ctx, kind, &entities, "-CreatedAt"); err != nil {
		return nil, errors.WithStack(err)
	}

	items := make([]*domain.Customer, 0, len(entities))
	for _, e := range entities {
		items = append(items, e.toDomain())
	}

	return items, nil
}

func (r *repository) GetAllWithPager(ctx context.Context, pager *domain.Pager) ([]*domain.Customer, error) {
	var entities []*entity

	if err := r.client.GetByFilter(
		ctx,
		kind,
		&entities,
		map[string]interface{}{},
		pager,
		"-CreatedAt"); err != nil {
		return nil, errors.WithStack(err)
	}

	items := make([]*domain.Customer, 0, len(entities))
	for _, e := range entities {
		items = append(items, e.toDomain())
	}

	return items, nil
}

func (r *repository) Get(ctx context.Context, id domain.CustomerID) (*domain.Customer, error) {
	entity := onlyID(id)

	if err := r.client.Get(ctx, entity); err != nil {
		return nil, errors.WithStack(err)
	}

	return entity.toDomain(), nil
}

func (r *repository) GetMulti(ctx context.Context, ids []domain.CustomerID) ([]*domain.Customer, error) {
	entities := make([]*entity, 0, len(ids))
	for _, id := range ids {
		entities = append(entities, onlyID(id))
	}

	if err := r.client.GetMulti(ctx, entities); err != nil {
		return nil, errors.WithStack(err)
	}

	items := make([]*domain.Customer, 0, len(entities))
	for _, e := range entities {
		items = append(items, e.toDomain())
	}

	return items, nil
}

func (r *repository) GetMultiWithIgnoreNotFound(ctx context.Context, ids []domain.CustomerID) ([]*domain.Customer, error) {
	entities := make([]*entity, 0, len(ids))
	for _, id := range ids {
		entities = append(entities, onlyID(id))
	}

	if err := r.client.GetMultiWithIgnoreError(ctx, entities); err != nil {
		return nil, errors.WithStack(err)
	}

	items := make([]*domain.Customer, 0, len(entities))
	for _, e := range entities {
		items = append(items, e.toDomain())
	}

	return items, nil
}

func (r *repository) GetTotalCount(ctx context.Context) (int64, error) {
	count, err := r.client.GetTotalCount(ctx, kind)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return count, nil
}

func (r *repository) Put(tx *boom.Transaction, customer *domain.Customer) error {
	if err := r.client.Put(tx, toEntity(customer)); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

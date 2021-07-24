package company_table

import (
	"context"

	"go.mercari.io/datastore/boom"

	"gae-go-sample/adapter"
	"gae-go-sample/domain"

	"github.com/pkg/errors"
)

func NewRepository(client adapter.DSClient) adapter.CompanyRepository {
	return &repository{
		client: client,
	}
}

type repository struct {
	client adapter.DSClient
}

func (r *repository) GetAll(ctx context.Context) ([]*domain.Company, error) {
	var entities []*entity

	if err := r.client.GetAll(ctx, kind, &entities, "-CreatedAt"); err != nil {
		return nil, errors.WithStack(err)
	}

	items := make([]*domain.Company, 0, len(entities))
	for _, e := range entities {
		items = append(items, e.toDomain())
	}

	return items, nil
}

func (r *repository) GetAllWithPager(ctx context.Context, pager *domain.Pager) ([]*domain.Company, error) {
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

	items := make([]*domain.Company, 0, len(entities))
	for _, e := range entities {
		items = append(items, e.toDomain())
	}

	return items, nil
}

func (r *repository) Get(ctx context.Context, id domain.CompanyID) (*domain.Company, error) {
	entity := onlyID(id)

	if err := r.client.Get(ctx, entity); err != nil {
		return nil, errors.WithStack(err)
	}

	return entity.toDomain(), nil
}

func (r *repository) GetMulti(ctx context.Context, ids []domain.CompanyID) ([]*domain.Company, error) {
	entities := make([]*entity, 0, len(ids))
	for _, id := range ids {
		entities = append(entities, onlyID(id))
	}

	if err := r.client.GetMulti(ctx, entities); err != nil {
		return nil, errors.WithStack(err)
	}

	items := make([]*domain.Company, 0, len(entities))
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

func (r *repository) Put(tx *boom.Transaction, company *domain.Company) error {
	if err := r.client.Put(tx, toEntity(company)); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

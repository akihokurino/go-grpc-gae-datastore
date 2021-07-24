package entry_table

import (
	"context"

	"go.mercari.io/datastore/boom"

	"gae-go-sample/adapter"
	"gae-go-sample/domain"

	"github.com/pkg/errors"
)

func NewRepository(client adapter.DSClient) adapter.EntryRepository {
	return &repository{
		client: client,
	}
}

type repository struct {
	client adapter.DSClient
}

func (r *repository) GetAllWithPager(ctx context.Context, pager *domain.Pager) ([]*domain.Entry, error) {
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

	items := make([]*domain.Entry, 0, len(entities))
	for _, e := range entities {
		items = append(items, e.toDomain())
	}

	return items, nil
}

func (r *repository) GetAllByProject(ctx context.Context, projectID domain.ProjectID) ([]*domain.Entry, error) {
	var entities []*entity

	if err := r.client.GetByFilter(
		ctx,
		kind,
		&entities,
		map[string]interface{}{
			"ProjectID =": projectID.String(),
		},
		nil,
		"-CreatedAt"); err != nil {
		return nil, errors.WithStack(err)
	}

	items := make([]*domain.Entry, 0, len(entities))
	for _, e := range entities {
		items = append(items, e.toDomain())
	}

	return items, nil
}

func (r *repository) GetAllByProjectWithPager(
	ctx context.Context,
	projectID domain.ProjectID,
	pager *domain.Pager) ([]*domain.Entry, error) {
	var entities []*entity

	if err := r.client.GetByFilter(
		ctx,
		kind,
		&entities,
		map[string]interface{}{
			"ProjectID =": projectID.String(),
		},
		pager,
		"-CreatedAt"); err != nil {
		return nil, errors.WithStack(err)
	}

	items := make([]*domain.Entry, 0, len(entities))
	for _, e := range entities {
		items = append(items, e.toDomain())
	}

	return items, nil
}

func (r *repository) GetAllByCustomerWithPager(
	ctx context.Context,
	customerID domain.CustomerID,
	pager *domain.Pager) ([]*domain.Entry, error) {
	var entities []*entity

	if err := r.client.GetByFilter(
		ctx,
		kind,
		&entities,
		map[string]interface{}{
			"CustomerID =": customerID.String(),
		},
		pager,
		"-CreatedAt"); err != nil {
		return nil, errors.WithStack(err)
	}

	items := make([]*domain.Entry, 0, len(entities))
	for _, e := range entities {
		items = append(items, e.toDomain())
	}

	return items, nil
}

func (r *repository) GetAllByCustomer(ctx context.Context, customerID domain.CustomerID) ([]*domain.Entry, error) {
	var entities []*entity

	if err := r.client.GetByFilter(
		ctx,
		kind,
		&entities,
		map[string]interface{}{
			"CustomerID =": customerID.String(),
		},
		nil,
		"-CreatedAt"); err != nil {
		return nil, errors.WithStack(err)
	}

	items := make([]*domain.Entry, 0, len(entities))
	for _, e := range entities {
		items = append(items, e.toDomain())
	}

	return items, nil
}

func (r *repository) Get(ctx context.Context, id domain.EntryID) (*domain.Entry, error) {
	entity := onlyID(id)

	if err := r.client.Get(ctx, entity); err != nil {
		return nil, errors.WithStack(err)
	}

	return entity.toDomain(), nil
}

func (r *repository) GetTotalCount(ctx context.Context) (int64, error) {
	count, err := r.client.GetTotalCount(ctx, kind)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return count, nil
}

func (r *repository) GetCountByProject(ctx context.Context, projectID domain.ProjectID) (int64, error) {
	count, err := r.client.GetCountByFilter(ctx, kind, map[string]interface{}{
		"ProjectID =": projectID.String(),
	})
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return count, nil
}

func (r *repository) GetCountByCustomer(ctx context.Context, customerID domain.CustomerID) (int64, error) {
	count, err := r.client.GetCountByFilter(ctx, kind, map[string]interface{}{
		"CustomerID =": customerID.String(),
	})
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return count, nil
}

func (r *repository) Exists(ctx context.Context, id domain.EntryID) (bool, error) {
	exists, err := r.client.Exists(ctx, onlyID(id))
	if err != nil {
		return false, errors.WithStack(err)
	}

	return exists, nil
}

func (r *repository) Put(tx *boom.Transaction, entry *domain.Entry) error {
	if err := r.client.Put(tx, toEntity(entry)); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *repository) Delete(tx *boom.Transaction, id domain.EntryID) error {
	if err := r.client.Delete(tx, onlyID(id)); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

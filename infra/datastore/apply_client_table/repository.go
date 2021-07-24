package apply_client_table

import (
	"context"

	"go.mercari.io/datastore/boom"

	"gae-go-sample/adapter"
	"gae-go-sample/domain"
	pb "gae-go-sample/proto/go/pb"

	"github.com/pkg/errors"
)

func NewRepository(client adapter.DSClient) adapter.ApplyClientRepository {
	return &repository{
		client: client,
	}
}

type repository struct {
	client adapter.DSClient
}

func (r *repository) GetAllWithPager(ctx context.Context, pager *domain.Pager) ([]*domain.ApplyClient, error) {
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

	items := make([]*domain.ApplyClient, 0, len(entities))
	for _, e := range entities {
		items = append(items, e.toDomain())
	}

	return items, nil
}

func (r *repository) GetAllByStatusWithPager(
	ctx context.Context,
	status pb.ApplyClient_Status,
	pager *domain.Pager) ([]*domain.ApplyClient, error) {
	var entities []*entity

	if err := r.client.GetByFilter(
		ctx,
		kind,
		&entities,
		map[string]interface{}{
			"Status =": int32(status),
		},
		pager,
		"-CreatedAt"); err != nil {
		return nil, errors.WithStack(err)
	}

	items := make([]*domain.ApplyClient, 0, len(entities))
	for _, e := range entities {
		items = append(items, e.toDomain())
	}

	return items, nil
}

func (r *repository) Get(ctx context.Context, email domain.ApplyClientID) (*domain.ApplyClient, error) {
	entity := onlyID(email)

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

func (r *repository) GetCountByStatus(ctx context.Context, status pb.ApplyClient_Status) (int64, error) {
	count, err := r.client.GetCountByFilter(ctx, kind, map[string]interface{}{
		"Status =": int32(status),
	})
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return count, nil
}

func (r *repository) Exists(ctx context.Context, email domain.ApplyClientID) (bool, error) {
	exists, err := r.client.Exists(ctx, onlyID(email))
	if err != nil {
		return false, errors.WithStack(err)
	}

	return exists, nil
}

func (r *repository) Put(tx *boom.Transaction, apply *domain.ApplyClient) error {
	if err := r.client.Put(tx, toEntity(apply)); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

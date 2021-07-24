package client_table

import (
	"context"

	"go.mercari.io/datastore/boom"

	"gae-go-sample/adapter"
	"gae-go-sample/domain"
	pb "gae-go-sample/proto/go/pb"

	"github.com/pkg/errors"
)

func NewRepository(client adapter.DSClient) adapter.ClientRepository {
	return &repository{
		client: client,
	}
}

type repository struct {
	client adapter.DSClient
}

func (r *repository) GetAll(ctx context.Context) ([]*domain.Client, error) {
	var entities []*entity

	if err := r.client.GetByFilter(
		ctx,
		kind,
		&entities,
		map[string]interface{}{
			"IsDeleted =": false,
		},
		nil,
		"-CreatedAt"); err != nil {
		return nil, errors.WithStack(err)
	}

	items := make([]*domain.Client, 0, len(entities))
	for _, e := range entities {
		items = append(items, e.toDomain())
	}

	return items, nil
}

func (r *repository) GetAllWithPager(ctx context.Context, pager *domain.Pager) ([]*domain.Client, error) {
	var entities []*entity

	if err := r.client.GetByFilter(
		ctx,
		kind,
		&entities,
		map[string]interface{}{
			"IsDeleted =": false,
		},
		pager,
		"-CreatedAt"); err != nil {
		return nil, errors.WithStack(err)
	}

	items := make([]*domain.Client, 0, len(entities))
	for _, e := range entities {
		items = append(items, e.toDomain())
	}

	return items, nil
}

func (r *repository) GetAllByCompany(
	ctx context.Context,
	companyID domain.CompanyID) ([]*domain.Client, error) {
	var entities []*entity

	if err := r.client.GetByFilter(
		ctx,
		kind,
		&entities,
		map[string]interface{}{
			"IsDeleted =": false,
			"CompanyID =": companyID.String(),
		},
		nil,
		"-CreatedAt"); err != nil {
		return nil, errors.WithStack(err)
	}

	items := make([]*domain.Client, 0, len(entities))
	for _, e := range entities {
		items = append(items, e.toDomain())
	}

	return items, nil
}

func (r *repository) GetAllByCompanyWithPager(
	ctx context.Context,
	companyID domain.CompanyID,
	pager *domain.Pager) ([]*domain.Client, error) {
	var entities []*entity

	if err := r.client.GetByFilter(
		ctx,
		kind,
		&entities,
		map[string]interface{}{
			"IsDeleted =": false,
			"CompanyID =": companyID.String(),
		},
		pager,
		"-CreatedAt"); err != nil {
		return nil, errors.WithStack(err)
	}

	items := make([]*domain.Client, 0, len(entities))
	for _, e := range entities {
		items = append(items, e.toDomain())
	}

	return items, nil
}

func (r *repository) Get(ctx context.Context, id domain.ClientID) (*domain.Client, error) {
	entity := onlyID(id)

	if err := r.client.Get(ctx, entity); err != nil {
		return nil, errors.WithStack(err)
	}

	if entity.IsDeleted {
		return nil, errors.WithStack(domain.ErrNoSuchEntity)
	}

	return entity.toDomain(), nil
}

func (r *repository) GetMulti(ctx context.Context, ids []domain.ClientID) ([]*domain.Client, error) {
	entities := make([]*entity, 0, len(ids))
	for _, id := range ids {
		entities = append(entities, onlyID(id))
	}

	if err := r.client.GetMulti(ctx, entities); err != nil {
		return nil, errors.WithStack(err)
	}

	items := make([]*domain.Client, 0, len(entities))
	for _, e := range entities {
		if e.IsDeleted {
			continue
		}
		items = append(items, e.toDomain())
	}

	return items, nil
}

func (r *repository) GetTotalCount(ctx context.Context) (int64, error) {
	count, err := r.client.GetCountByFilter(ctx, kind, map[string]interface{}{
		"IsDeleted =": false,
	})
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return count, nil
}

func (r *repository) GetCountByCompany(ctx context.Context, companyID domain.CompanyID) (int64, error) {
	count, err := r.client.GetCountByFilter(ctx, kind, map[string]interface{}{
		"IsDeleted =": false,
		"CompanyID =": companyID.String(),
	})
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return count, nil
}

func (r *repository) GetCountByCompanyAndRole(ctx context.Context, companyID domain.CompanyID, role pb.Client_Role) (int64, error) {
	count, err := r.client.GetCountByFilter(ctx, kind, map[string]interface{}{
		"IsDeleted =": false,
		"CompanyID =": companyID.String(),
		"Role =":      int32(role),
	})
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return count, nil
}

func (r *repository) Put(tx *boom.Transaction, client *domain.Client) error {
	if err := r.client.Put(tx, toEntity(client)); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

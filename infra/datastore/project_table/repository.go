package project_table

import (
	"context"

	"go.mercari.io/datastore/boom"

	"gae-go-recruiting-server/adapter"
	"gae-go-recruiting-server/domain"
	pb "gae-go-recruiting-server/proto/go/pb"

	"github.com/pkg/errors"
)

func NewRepository(client adapter.DSClient) adapter.ProjectRepository {
	return &repository{
		client: client,
	}
}

type repository struct {
	client adapter.DSClient
}

func (r *repository) GetAll(ctx context.Context) ([]*domain.Project, error) {
	var entities []*entity

	if err := r.client.GetAll(ctx, kind, &entities, "-CreatedAt"); err != nil {
		return nil, errors.WithStack(err)
	}

	items := make([]*domain.Project, 0, len(entities))
	for _, e := range entities {
		items = append(items, e.toDomain())
	}

	return items, nil
}

func (r *repository) GetAllWithPager(ctx context.Context, pager *domain.Pager) ([]*domain.Project, error) {
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

	items := make([]*domain.Project, 0, len(entities))
	for _, e := range entities {
		items = append(items, e.toDomain())
	}

	return items, nil
}

func (r *repository) GetAllByCompany(ctx context.Context, companyID domain.CompanyID) ([]*domain.Project, error) {
	var entities []*entity

	if err := r.client.GetByFilter(
		ctx,
		kind,
		&entities,
		map[string]interface{}{
			"CompanyID =": companyID.String(),
		},
		nil,
		"-CreatedAt"); err != nil {
		return nil, errors.WithStack(err)
	}

	items := make([]*domain.Project, 0, len(entities))
	for _, e := range entities {
		items = append(items, e.toDomain())
	}

	return items, nil
}

func (r *repository) GetAllByCompanyAndStatus(
	ctx context.Context,
	companyID domain.CompanyID,
	status pb.Project_Status) ([]*domain.Project, error) {
	var entities []*entity

	if err := r.client.GetByFilter(
		ctx,
		kind,
		&entities,
		map[string]interface{}{
			"CompanyID =": companyID.String(),
			"Status =":    int32(status),
		},
		nil,
		"-CreatedAt"); err != nil {
		return nil, errors.WithStack(err)
	}

	items := make([]*domain.Project, 0, len(entities))
	for _, e := range entities {
		items = append(items, e.toDomain())
	}

	return items, nil
}

func (r *repository) GetAllByCompanyWithPager(
	ctx context.Context,
	companyID domain.CompanyID,
	pager *domain.Pager) ([]*domain.Project, error) {
	var entities []*entity

	if err := r.client.GetByFilter(
		ctx,
		kind,
		&entities,
		map[string]interface{}{
			"CompanyID =": companyID.String(),
		},
		pager,
		"-CreatedAt"); err != nil {
		return nil, errors.WithStack(err)
	}

	items := make([]*domain.Project, 0, len(entities))
	for _, e := range entities {
		items = append(items, e.toDomain())
	}

	return items, nil
}

func (r *repository) GetAllByStatusWithPager(
	ctx context.Context,
	status pb.Project_Status,
	pager *domain.Pager) ([]*domain.Project, error) {
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

	items := make([]*domain.Project, 0, len(entities))
	for _, e := range entities {
		items = append(items, e.toDomain())
	}

	return items, nil
}

func (r *repository) GetAllByStatusAndEntryNumOrderWithPager(
	ctx context.Context,
	status pb.Project_Status,
	pager *domain.Pager) ([]*domain.Project, error) {
	var entities []*entity

	if err := r.client.GetByFilter(
		ctx,
		kind,
		&entities,
		map[string]interface{}{
			"Status =": int32(status),
		},
		pager,
		"-EntryNum"); err != nil {
		return nil, errors.WithStack(err)
	}

	items := make([]*domain.Project, 0, len(entities))
	for _, e := range entities {
		items = append(items, e.toDomain())
	}

	return items, nil
}

func (r *repository) Get(ctx context.Context, id domain.ProjectID) (*domain.Project, error) {
	entity := onlyID(id)

	if err := r.client.Get(ctx, entity); err != nil {
		return nil, errors.WithStack(err)
	}

	return entity.toDomain(), nil
}

func (r *repository) GetMulti(ctx context.Context, ids []domain.ProjectID) ([]*domain.Project, error) {
	entities := make([]*entity, 0, len(ids))
	for _, id := range ids {
		entities = append(entities, onlyID(id))
	}

	if err := r.client.GetMulti(ctx, entities); err != nil {
		return nil, errors.WithStack(err)
	}

	items := make([]*domain.Project, 0, len(entities))
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

func (r *repository) GetCountByCompany(ctx context.Context, companyID domain.CompanyID) (int64, error) {
	count, err := r.client.GetCountByFilter(ctx, kind, map[string]interface{}{
		"CompanyID =": companyID.String(),
	})
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return count, nil
}

func (r *repository) GetCountByStatus(ctx context.Context, status pb.Project_Status) (int64, error) {
	count, err := r.client.GetCountByFilter(ctx, kind, map[string]interface{}{
		"Status =": int32(status),
	})
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return count, nil
}

func (r *repository) GetCountByCompanyAndStatus(ctx context.Context, companyID domain.CompanyID, status pb.Project_Status) (int64, error) {
	count, err := r.client.GetCountByFilter(ctx, kind, map[string]interface{}{
		"CompanyID =": companyID.String(),
		"Status =":    int32(status),
	})
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return count, nil
}

func (r *repository) Put(tx *boom.Transaction, project *domain.Project) error {
	if err := r.client.Put(tx, toEntity(project)); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *repository) Delete(tx *boom.Transaction, id domain.ProjectID) error {
	if err := r.client.Delete(tx, onlyID(id)); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

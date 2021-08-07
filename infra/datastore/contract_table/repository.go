package contract_table

import (
	"context"

	"go.mercari.io/datastore/boom"

	"gae-go-recruiting-server/adapter"
	"gae-go-recruiting-server/domain"

	"github.com/pkg/errors"
)

func NewRepository(client adapter.DSClient) adapter.ContractRepository {
	return &repository{
		client: client,
	}
}

type repository struct {
	client adapter.DSClient
}

func (r *repository) GetAllWithPager(ctx context.Context, pager *domain.Pager) ([]*domain.Contract, error) {
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

	items := make([]*domain.Contract, 0, len(entities))
	for _, e := range entities {
		item, err := e.toDomain()
		if err != nil {
			return nil, errors.WithStack(err)
		}
		items = append(items, item)
	}

	return items, nil
}

func (r *repository) GetAllByCompanyWithPager(
	ctx context.Context,
	companyID domain.CompanyID,
	pager *domain.Pager) ([]*domain.Contract, error) {
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

	items := make([]*domain.Contract, 0, len(entities))
	for _, e := range entities {
		item, err := e.toDomain()
		if err != nil {
			return nil, errors.WithStack(err)
		}
		items = append(items, item)
	}

	return items, nil
}

func (r *repository) Get(ctx context.Context, id domain.ContractID) (*domain.Contract, error) {
	entity := onlyID(id)

	if err := r.client.Get(ctx, entity); err != nil {
		return nil, errors.WithStack(err)
	}

	item, err := entity.toDomain()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return item, nil
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

func (r *repository) Exists(ctx context.Context, id domain.ContractID) (bool, error) {
	exists, err := r.client.Exists(ctx, onlyID(id))
	if err != nil {
		return false, errors.WithStack(err)
	}

	return exists, nil
}

func (r *repository) Put(tx *boom.Transaction, contract *domain.Contract) error {
	if err := r.client.Put(tx, toEntity(contract)); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *repository) Delete(tx *boom.Transaction, id domain.ContractID) error {
	if err := r.client.Delete(tx, onlyID(id)); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

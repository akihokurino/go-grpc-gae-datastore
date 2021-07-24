package message_room_table

import (
	"context"

	"go.mercari.io/datastore/boom"

	"gae-go-sample/adapter"
	"gae-go-sample/domain"

	"github.com/pkg/errors"
)

func NewRepository(client adapter.DSClient) adapter.MessageRoomRepository {
	return &repository{
		client: client,
	}
}

type repository struct {
	client adapter.DSClient
}

func (r *repository) GetAll(ctx context.Context) ([]*domain.MessageRoom, error) {
	var entities []*entity

	if err := r.client.GetAll(ctx, kind, &entities, "-UpdatedAt"); err != nil {
		return nil, errors.WithStack(err)
	}

	items := make([]*domain.MessageRoom, 0, len(entities))
	for _, e := range entities {
		items = append(items, e.toDomain())
	}

	return items, nil
}

func (r *repository) GetAllByCustomer(ctx context.Context, customerID domain.CustomerID) ([]*domain.MessageRoom, error) {
	var entities []*entity

	if err := r.client.GetByFilter(
		ctx,
		kind,
		&entities,
		map[string]interface{}{
			"CustomerID =": customerID.String(),
		},
		nil,
		"-UpdatedAt"); err != nil {
		return nil, errors.WithStack(err)
	}

	items := make([]*domain.MessageRoom, 0, len(entities))
	for _, e := range entities {
		items = append(items, e.toDomain())
	}

	return items, nil
}

func (r *repository) GetAllByCustomerWithPager(
	ctx context.Context,
	customerID domain.CustomerID,
	pager *domain.Pager) ([]*domain.MessageRoom, error) {
	var entities []*entity

	if err := r.client.GetByFilter(
		ctx,
		kind,
		&entities,
		map[string]interface{}{
			"CustomerID =": customerID.String(),
		},
		pager,
		"-UpdatedAt"); err != nil {
		return nil, errors.WithStack(err)
	}

	items := make([]*domain.MessageRoom, 0, len(entities))
	for _, e := range entities {
		items = append(items, e.toDomain())
	}

	return items, nil
}

func (r *repository) GetAllByCompany(ctx context.Context, companyID domain.CompanyID) ([]*domain.MessageRoom, error) {
	var entities []*entity

	if err := r.client.GetByFilter(
		ctx,
		kind,
		&entities,
		map[string]interface{}{
			"CompanyID =": companyID.String(),
		},
		nil,
		"-UpdatedAt"); err != nil {
		return nil, errors.WithStack(err)
	}

	items := make([]*domain.MessageRoom, 0, len(entities))
	for _, e := range entities {
		items = append(items, e.toDomain())
	}

	return items, nil
}

func (r *repository) GetAllByCompanyWithPager(
	ctx context.Context,
	companyID domain.CompanyID,
	pager *domain.Pager) ([]*domain.MessageRoom, error) {
	var entities []*entity

	if err := r.client.GetByFilter(
		ctx,
		kind,
		&entities,
		map[string]interface{}{
			"CompanyID =": companyID.String(),
		},
		pager,
		"-UpdatedAt"); err != nil {
		return nil, errors.WithStack(err)
	}

	items := make([]*domain.MessageRoom, 0, len(entities))
	for _, e := range entities {
		items = append(items, e.toDomain())
	}

	return items, nil
}

func (r *repository) GetLastByProjectAndCustomer(
	ctx context.Context,
	projectID domain.ProjectID,
	customerID domain.CustomerID) (*domain.MessageRoom, error) {
	var entities []*entity

	if err := r.client.GetLastByFilter(
		ctx,
		kind,
		&entities,
		map[string]interface{}{
			"ProjectID =":  projectID.String(),
			"CustomerID =": customerID.String(),
		},
		"-UpdatedAt"); err != nil {
		return nil, errors.WithStack(err)
	}

	return entities[0].toDomain(), nil
}

func (r *repository) Get(ctx context.Context, id domain.MessageRoomID) (*domain.MessageRoom, error) {
	entity := onlyID(id)

	if err := r.client.Get(ctx, entity); err != nil {
		return nil, errors.WithStack(err)
	}

	return entity.toDomain(), nil
}

func (r *repository) Exists(ctx context.Context, id domain.MessageRoomID) (bool, error) {
	exists, err := r.client.Exists(ctx, onlyID(id))
	if err != nil {
		return false, errors.WithStack(err)
	}

	return exists, nil
}

func (r *repository) Put(tx *boom.Transaction, room *domain.MessageRoom) error {
	if err := r.client.Put(tx, toEntity(room)); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *repository) Delete(tx *boom.Transaction, id domain.MessageRoomID) error {
	if err := r.client.Delete(tx, onlyID(id)); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

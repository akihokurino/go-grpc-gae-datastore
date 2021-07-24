package user_table

import (
	"context"

	"go.mercari.io/datastore/boom"

	"gae-go-sample/adapter"
	"gae-go-sample/domain"

	"github.com/pkg/errors"
)

func NewRepository(client adapter.DSClient) adapter.UserRepository {
	return &repository{
		client: client,
	}
}

type repository struct {
	client adapter.DSClient
}

func (r *repository) Get(ctx context.Context, email domain.UserID) (*domain.User, error) {
	entity := onlyID(email)

	if err := r.client.Get(ctx, entity); err != nil {
		return nil, errors.WithStack(err)
	}

	return entity.toDomain(), nil
}

func (r *repository) Exists(ctx context.Context, email domain.UserID) (bool, error) {
	exists, err := r.client.Exists(ctx, onlyID(email))
	if err != nil {
		return false, errors.WithStack(err)
	}

	return exists, nil
}

func (r *repository) Put(tx *boom.Transaction, user *domain.User) error {
	if err := r.client.Put(tx, toEntity(user)); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

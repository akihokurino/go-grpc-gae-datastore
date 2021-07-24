package message_table

import (
	"context"

	"go.mercari.io/datastore/boom"

	"gae-go-sample/adapter"
	"gae-go-sample/domain"

	"github.com/pkg/errors"
)

func NewRepository(client adapter.DSClient) adapter.MessageRepository {
	return &repository{
		client: client,
	}
}

type repository struct {
	client adapter.DSClient
}

func (r *repository) GetAllByRoom(ctx context.Context, roomID domain.MessageRoomID) ([]*domain.Message, error) {
	var entities []*entity

	if err := r.client.GetByFilter(
		ctx,
		kind,
		&entities,
		map[string]interface{}{
			"RoomID =": roomID.String(),
		},
		nil,
		"-CreatedAt"); err != nil {
		return nil, errors.WithStack(err)
	}

	items := make([]*domain.Message, 0, len(entities))
	for _, e := range entities {
		items = append(items, e.toDomain())
	}

	return items, nil
}

func (r *repository) GetAllByRoomWithPager(
	ctx context.Context,
	roomID domain.MessageRoomID,
	pager *domain.Pager) ([]*domain.Message, error) {
	var entities []*entity

	if err := r.client.GetByFilter(
		ctx,
		kind,
		&entities,
		map[string]interface{}{
			"RoomID =": roomID.String(),
		},
		pager,
		"-CreatedAt"); err != nil {
		return nil, errors.WithStack(err)
	}
	items := make([]*domain.Message, 0, len(entities))
	for _, e := range entities {
		items = append(items, e.toDomain())
	}

	return items, nil
}

func (r *repository) GetLastByRoom(ctx context.Context, roomID domain.MessageRoomID) (*domain.Message, error) {
	var entities []*entity

	if err := r.client.GetLastByFilter(
		ctx,
		kind,
		&entities,
		map[string]interface{}{
			"RoomID =": roomID.String(),
		},
		"-CreatedAt"); err != nil {
		return nil, errors.WithStack(err)
	}

	return entities[0].toDomain(), nil
}

func (r *repository) Get(ctx context.Context, id domain.MessageID) (*domain.Message, error) {
	entity := onlyID(id)

	if err := r.client.Get(ctx, entity); err != nil {
		return nil, errors.WithStack(err)
	}

	return entity.toDomain(), nil
}

func (r *repository) Put(tx *boom.Transaction, message *domain.Message) error {
	if err := r.client.Put(tx, toEntity(message)); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *repository) Delete(tx *boom.Transaction, id domain.MessageID) error {
	if err := r.client.Delete(tx, onlyID(id)); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

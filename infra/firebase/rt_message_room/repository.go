package rt_message_room

import (
	"context"

	"gae-go-sample/adapter"
	"gae-go-sample/domain"

	"github.com/pkg/errors"
)

type rtMessageRoomRepository struct {
	fireClient adapter.FirebaseClient
}

func NewRtMessageRoomRepository(fireClient adapter.FirebaseClient) adapter.RtMessageRoomRepository {
	return &rtMessageRoomRepository{
		fireClient: fireClient,
	}
}

func (r *rtMessageRoomRepository) Put(
	ctx context.Context,
	room *domain.MessageRoom,
	customer *domain.Customer,
	clients []*domain.Client) error {
	rtdbClient, err := r.fireClient.RTDBClient(ctx)
	if err != nil {
		return errors.WithStack(err)
	}

	roomEntity := newMessageRoom()
	roomRef := rtdbClient.NewRef(messageRoomPath(room.ID()))

	if err := roomRef.Set(ctx, roomEntity); err != nil {
		return errors.WithStack(err)
	}

	memberEntity := newMemberFromCustomer(customer)
	memberRef := rtdbClient.NewRef(memberPath(room.ID(), domain.UserID(customer.ID)))

	if err := memberRef.Set(ctx, memberEntity); err != nil {
		return errors.WithStack(err)
	}

	for _, client := range clients {
		memberEntity = newMemberFromClient(client)
		memberRef = rtdbClient.NewRef(memberPath(room.ID(), domain.UserID(client.ID)))

		if err := memberRef.Set(ctx, memberEntity); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func (r *rtMessageRoomRepository) Delete(
	ctx context.Context,
	room *domain.MessageRoom) error {
	rtdbClient, err := r.fireClient.RTDBClient(ctx)
	if err != nil {
		return errors.WithStack(err)
	}

	roomRef := rtdbClient.NewRef(messageRoomPath(room.ID()))

	if err := roomRef.Delete(ctx); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

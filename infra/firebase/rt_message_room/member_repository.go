package rt_message_room

import (
	"context"

	"gae-go-recruiting-server/adapter"
	"gae-go-recruiting-server/domain"

	"github.com/pkg/errors"
)

type rtMemberRepository struct {
	fireClient adapter.FirebaseClient
}

func NewRtMemberRepository(fireClient adapter.FirebaseClient) adapter.RtMemberRepository {
	return &rtMemberRepository{
		fireClient: fireClient,
	}
}

func (r *rtMemberRepository) UpdateCustomer(
	ctx context.Context,
	room *domain.MessageRoom,
	customer *domain.Customer) error {
	rtdbClient, err := r.fireClient.RTDBClient(ctx)
	if err != nil {
		return errors.WithStack(err)
	}

	memberEntity := newMemberFromCustomer(customer)
	memberRef := rtdbClient.NewRef(memberPath(room.ID(), domain.UserID(customer.ID)))

	if err := memberRef.Set(ctx, memberEntity); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *rtMemberRepository) UpdateClient(
	ctx context.Context,
	room *domain.MessageRoom,
	client *domain.Client) error {
	rtdbClient, err := r.fireClient.RTDBClient(ctx)
	if err != nil {
		return errors.WithStack(err)
	}

	memberEntity := newMemberFromClient(client)
	memberRef := rtdbClient.NewRef(memberPath(room.ID(), domain.UserID(client.ID)))

	if err := memberRef.Set(ctx, memberEntity); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

package message_room_app

import (
	"context"

	"go.mercari.io/datastore/boom"

	"gae-go-sample/domain"

	"github.com/pkg/errors"
)

type clientApplication struct {
	executorID domain.ClientID
	*application
}

func (a *clientApplication) GetAllWithPager(
	ctx context.Context,
	page int32,
	offset int32) ([]*domain.MessageRoom, error) {
	me, err := a.clientRepository.Get(ctx, a.executorID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	pager := domain.NewPager(page, offset)

	rooms, err := a.messageRoomRepository.GetAllByCompanyWithPager(ctx, me.CompanyID, pager)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return rooms, nil
}

func (a *clientApplication) Get(
	ctx context.Context,
	projectID domain.ProjectID,
	customerID domain.CustomerID) (*domain.MessageRoom, error) {
	me, err := a.clientRepository.Get(ctx, a.executorID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	room, err := a.messageRoomRepository.Get(ctx, domain.NewMessageRoomID(projectID, customerID, me.CompanyID))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return room, nil
}

func (a *clientApplication) Read(
	ctx context.Context,
	projectID domain.ProjectID,
	customerID domain.CustomerID) error {
	client, err := a.clientRepository.Get(ctx, a.executorID)
	if err != nil {
		return errors.WithStack(err)
	}

	room, err := a.messageRoomRepository.Get(ctx, domain.NewMessageRoomID(projectID, customerID, client.CompanyID))
	if err != nil {
		return errors.WithStack(err)
	}

	me, err := a.userRepository.Get(ctx, domain.UserID(a.executorID))
	if err != nil {
		return errors.WithStack(err)
	}

	if !me.IsClient() {
		return errors.WithStack(domain.ErrInvalidUserRole)
	}

	isUnRead, err := room.IsUnRead(domain.ByCompany)
	if err != nil {
		return errors.WithStack(err)
	}

	if !isUnRead {
		return nil
	}

	if err := room.Read(domain.ByCompany); err != nil {
		return errors.WithStack(err)
	}

	if err := a.transaction(ctx, func(tx *boom.Transaction) error {
		if err := a.messageRoomRepository.Put(tx, room); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (a *clientApplication) Delete(
	ctx context.Context,
	projectID domain.ProjectID,
	customerID domain.CustomerID) error {
	me, err := a.clientRepository.Get(ctx, a.executorID)
	if err != nil {
		return errors.WithStack(err)
	}

	room, err := a.messageRoomRepository.Get(ctx, domain.NewMessageRoomID(projectID, customerID, me.CompanyID))
	if err != nil {
		return errors.WithStack(err)
	}

	if !room.IsIncludeClient(me) {
		return errors.WithStack(domain.ErrUserIsNotMember)
	}

	messages, err := a.messageRepository.GetAllByRoom(ctx, room.ID())

	if err := a.transaction(ctx, func(tx *boom.Transaction) error {
		for _, message := range messages {
			if err := a.messageRepository.Delete(tx, message.ID); err != nil {
				return err
			}
		}

		if err := a.messageRoomRepository.Delete(tx, room.ID()); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return errors.WithStack(err)
	}

	if err := a.rtMessageRoomRepository.Delete(ctx, room); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

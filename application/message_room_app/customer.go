package message_room_app

import (
	"context"

	"go.mercari.io/datastore/boom"

	"gae-go-recruiting-server/domain"

	"github.com/pkg/errors"
)

type customerApplication struct {
	executorID domain.CustomerID
	*application
}

func (a *customerApplication) GetAllWithPager(
	ctx context.Context,
	page int32,
	offset int32) ([]*domain.MessageRoom, error) {
	pager := domain.NewPager(page, offset)

	rooms, err := a.messageRoomRepository.GetAllByCustomerWithPager(ctx, a.executorID, pager)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return rooms, nil
}

func (a *customerApplication) Get(
	ctx context.Context,
	projectID domain.ProjectID,
	companyID domain.CompanyID) (*domain.MessageRoom, error) {
	room, err := a.messageRoomRepository.Get(ctx, domain.NewMessageRoomID(projectID, a.executorID, companyID))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return room, nil
}

func (a *customerApplication) Read(
	ctx context.Context,
	projectID domain.ProjectID,
	companyID domain.CompanyID) error {
	room, err := a.messageRoomRepository.Get(ctx, domain.NewMessageRoomID(projectID, a.executorID, companyID))
	if err != nil {
		return errors.WithStack(err)
	}

	me, err := a.userRepository.Get(ctx, domain.UserID(a.executorID))
	if err != nil {
		return errors.WithStack(err)
	}

	if !me.IsCustomer() {
		return errors.WithStack(domain.ErrInvalidUserRole)
	}

	isUnRead, err := room.IsUnRead(domain.ByCustomer)
	if err != nil {
		return errors.WithStack(err)
	}

	if !isUnRead {
		return nil
	}

	if err := room.Read(domain.ByCustomer); err != nil {
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

func (a *customerApplication) Delete(
	ctx context.Context,
	projectID domain.ProjectID,
	companyID domain.CompanyID) error {
	me, err := a.customerRepository.Get(ctx, a.executorID)
	if err != nil {
		return errors.WithStack(err)
	}

	room, err := a.messageRoomRepository.Get(ctx, domain.NewMessageRoomID(projectID, a.executorID, companyID))
	if err != nil {
		return errors.WithStack(err)
	}

	if !room.IsIncludeCustomer(me) {
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

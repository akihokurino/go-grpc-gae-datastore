package message_room_app

import (
	"context"
	"sync"

	"gae-go-recruiting-server/domain"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

type adminApplication struct {
	executorID domain.AdminUserID
	*application
}

func (a *adminApplication) GetAllByIDs(
	ctx context.Context,
	ids []struct {
		ProjectID  domain.ProjectID
		CustomerID domain.CustomerID
	}) ([]*domain.MessageRoom, error) {
	rooms := make([]*domain.MessageRoom, 0, len(ids))

	lock := new(sync.Mutex)
	eg := errgroup.Group{}

	for _, id := range ids {
		roomID := id

		eg.Go(func() error {
			room, err := a.messageRoomRepository.GetLastByProjectAndCustomer(ctx, roomID.ProjectID, roomID.CustomerID)
			if err != nil {
				return err
			}

			lock.Lock()
			rooms = append(rooms, room)
			lock.Unlock()

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, errors.WithStack(err)
	}

	return rooms, nil
}

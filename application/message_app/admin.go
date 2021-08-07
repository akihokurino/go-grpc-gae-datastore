package message_app

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

func (a *adminApplication) GetAllByRoomWithPager(
	ctx context.Context,
	roomID domain.MessageRoomID,
	page int32,
	offset int32) ([]*domain.Message, error) {
	pager := domain.NewPager(page, offset)

	messages, err := a.messageRepository.GetAllByRoomWithPager(
		ctx,
		roomID,
		pager)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for i := range messages {
		imageURLWithSignature, err := a.publishResourceService(ctx, messages[i].GSImageURL)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		messages[i].SignedImageURL = imageURLWithSignature

		fileURLWithSignature, err := a.publishResourceService(ctx, messages[i].GSFileURL)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		messages[i].SignedFileURL = fileURLWithSignature
	}

	return messages, nil
}

func (a *adminApplication) GetAllNewestByRooms(
	ctx context.Context,
	roomIDs []domain.MessageRoomID) ([]*domain.Message, error) {
	messageMap := make(map[domain.MessageRoomID]*domain.Message, 0)

	mutex := sync.Mutex{}
	eg := errgroup.Group{}

	for i := range roomIDs {
		id := roomIDs[i]

		eg.Go(func() error {
			message, err := a.messageRepository.GetLastByRoom(ctx, id)
			if err != nil && !domain.IsNoSuchEntityErr(err) {
				return err
			}
			if err != nil && domain.IsNoSuchEntityErr(err) {
				return nil
			}

			mutex.Lock()
			messageMap[id] = message
			mutex.Unlock()

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, errors.WithStack(err)
	}

	messages := make([]*domain.Message, 0, len(messageMap))
	for _, message := range messageMap {
		messages = append(messages, message)
	}

	for i := range messages {
		imageURLWithSignature, err := a.publishResourceService(ctx, messages[i].GSImageURL)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		messages[i].SignedImageURL = imageURLWithSignature

		fileURLWithSignature, err := a.publishResourceService(ctx, messages[i].GSFileURL)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		messages[i].SignedFileURL = fileURLWithSignature
	}

	return messages, nil
}

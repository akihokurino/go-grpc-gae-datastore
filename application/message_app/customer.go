package message_app

import (
	"context"
	"sync"

	"gae-go-sample/domain"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

type customerApplication struct {
	executorID domain.CustomerID
	*application
}

func (a *customerApplication) GetAllByRoomWithPager(
	ctx context.Context,
	projectID domain.ProjectID,
	companyID domain.CompanyID,
	page int32,
	offset int32) ([]*domain.Message, error) {
	pager := domain.NewPager(page, offset)

	messages, err := a.messageRepository.GetAllByRoomWithPager(
		ctx,
		domain.NewMessageRoomID(projectID, a.executorID, companyID),
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

func (a *customerApplication) GetAllNewestByRooms(
	ctx context.Context,
	roomIDParams []struct {
		ProjectID domain.ProjectID
		CompanyID domain.CompanyID
	}) ([]*domain.Message, error) {
	roomIDs := make([]domain.MessageRoomID, 0, len(roomIDParams))
	for _, id := range roomIDParams {
		roomIDs = append(roomIDs, domain.NewMessageRoomID(id.ProjectID, a.executorID, id.CompanyID))
	}

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

package entry_app

import (
	"context"

	"golang.org/x/sync/errgroup"

	"gae-go-recruiting-server/domain"

	"github.com/pkg/errors"
)

type clientApplication struct {
	executorID domain.ClientID
	*application
}

func (a *clientApplication) GetAllByProject(
	ctx context.Context,
	projectID domain.ProjectID) ([]*domain.Entry, error) {
	me, err := a.clientRepository.Get(ctx, a.executorID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if _, err := a.validProjectService(ctx, me, projectID); err != nil {
		return nil, errors.WithStack(err)
	}

	entries, err := a.entryRepository.GetAllByProject(ctx, projectID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return entries, nil
}

func (a *clientApplication) GetCountByProjects(
	ctx context.Context,
	projectIDs []domain.ProjectID) (map[domain.ProjectID]int64, error) {
	me, err := a.clientRepository.Get(ctx, a.executorID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	countMap := make(map[domain.ProjectID]int64, 0)

	eg := errgroup.Group{}

	for _, id := range projectIDs {
		projectID := id

		eg.Go(func() error {
			if _, err := a.validProjectService(ctx, me, projectID); err != nil {
				return err
			}

			count, err := a.entryRepository.GetCountByProject(ctx, projectID)
			if err != nil {
				return err
			}

			countMap[projectID] = count

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, errors.WithStack(err)
	}

	return countMap, nil
}

func (a *clientApplication) GetCountByProject(
	ctx context.Context,
	projectID domain.ProjectID) (int64, error) {
	me, err := a.clientRepository.Get(ctx, a.executorID)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	if _, err := a.validProjectService(ctx, me, projectID); err != nil {
		return 0, errors.WithStack(err)
	}

	count, err := a.entryRepository.GetCountByProject(ctx, projectID)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return count, nil
}

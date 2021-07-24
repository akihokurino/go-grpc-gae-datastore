package entry_app

import (
	"context"

	"gae-go-sample/domain"

	"github.com/pkg/errors"
)

type adminApplication struct {
	executorID domain.AdminUserID
	*application
}

func (a *adminApplication) GetAllWithPager(ctx context.Context, page int32, offset int32) ([]*domain.Entry, error) {
	entries, err := a.entryRepository.GetAllWithPager(ctx, domain.NewPager(page, offset))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return entries, nil
}

func (a *adminApplication) GetAllByProjectWithPager(
	ctx context.Context,
	projectID domain.ProjectID,
	page int32,
	offset int32) ([]*domain.Entry, error) {
	entries, err := a.entryRepository.GetAllByProjectWithPager(ctx, projectID, domain.NewPager(page, offset))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return entries, nil
}

func (a *adminApplication) Get(ctx context.Context, customerID domain.CustomerID, projectID domain.ProjectID) (*domain.Entry, error) {
	entry, err := a.entryRepository.Get(ctx, domain.NewEntryID(customerID, projectID))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return entry, nil
}

func (a *adminApplication) GetTotalCount(ctx context.Context) (int64, error) {
	count, err := a.entryRepository.GetTotalCount(ctx)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return count, nil
}

func (a *adminApplication) GetCountByProject(ctx context.Context, projectID domain.ProjectID) (int64, error) {
	count, err := a.entryRepository.GetCountByProject(ctx, projectID)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return count, nil
}

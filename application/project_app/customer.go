package project_app

import (
	"context"

	"gae-go-recruiting-server/adapter"
	"gae-go-recruiting-server/domain"
	pb "gae-go-recruiting-server/proto/go/pb"

	"github.com/pkg/errors"
)

type customerApplication struct {
	executorID domain.CustomerID
	*application
}

func (a *customerApplication) GetAllByOpenWithPager(ctx context.Context, page int32, offset int32) ([]*domain.Project, error) {
	projects, err := a.projectRepository.GetAllByStatusWithPager(ctx, pb.Project_Status_Open, domain.NewPager(page, offset))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for i := range projects {
		urlWithSignature, err := a.publishResourceService(ctx, projects[i].GSThumbnailURL)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		projects[i].SignedThumbnailURL = urlWithSignature
	}

	return projects, nil
}

// TODO: おすすめ一覧を表示する
func (a *customerApplication) GetAllByRecommend(ctx context.Context) ([]*domain.Project, error) {
	projects, err := a.projectRepository.GetAllByStatusWithPager(ctx, pb.Project_Status_Open, domain.NewPager(1, 10))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for i := range projects {
		urlWithSignature, err := a.publishResourceService(ctx, projects[i].GSThumbnailURL)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		projects[i].SignedThumbnailURL = urlWithSignature
	}

	return projects, nil
}

func (a *customerApplication) GetAllByEntryNumOrderWithPager(
	ctx context.Context,
	page int32,
	offset int32) ([]*domain.Project, error) {
	pager := domain.NewPager(page, offset)

	projects, err := a.projectRepository.GetAllByStatusAndEntryNumOrderWithPager(
		ctx,
		pb.Project_Status_Open,
		pager)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for i := range projects {
		urlWithSignature, err := a.publishResourceService(ctx, projects[i].GSThumbnailURL)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		projects[i].SignedThumbnailURL = urlWithSignature
	}

	return projects, nil
}

func (a *customerApplication) SearchWithPager(
	ctx context.Context,
	params adapter.SearchProjectParams,
	page int32,
	offset int32) ([]*domain.Project, error) {
	pager := domain.NewPager(page, offset)

	projectIDsWithHighlights, err := a.projectIndexRepository.SearchByStatusWithConditionWithPager(
		ctx,
		params.Query,
		pb.Project_Status_Open,
		pager)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	ids := make([]domain.ProjectID, 0, len(projectIDsWithHighlights))
	for _, id := range projectIDsWithHighlights {
		ids = append(ids, id.ID)
	}

	projects, err := a.projectRepository.GetMulti(ctx, ids)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for i := range projects {
		projectID := projects[i].ID
		for _, idWithHighlights := range projectIDsWithHighlights {
			if projectID == idWithHighlights.ID {
				projects[i].Highlights = idWithHighlights.Highlights
				break
			}
		}
	}

	for i := range projects {
		urlWithSignature, err := a.publishResourceService(ctx, projects[i].GSThumbnailURL)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		projects[i].SignedThumbnailURL = urlWithSignature
	}

	return projects, nil
}

func (a *customerApplication) GetAllByEntryWithPager(
	ctx context.Context,
	page int32,
	offset int32) ([]*domain.Project, error) {
	me, err := a.customerRepository.Get(ctx, a.executorID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	entries, err := a.entryRepository.GetAllByCustomerWithPager(ctx, me.ID, domain.NewPager(page, offset))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	projectIDs := make([]domain.ProjectID, 0, len(entries))
	for _, entry := range entries {
		projectIDs = append(projectIDs, entry.ProjectID)
	}

	projects, err := a.projectRepository.GetMulti(ctx, projectIDs)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for i := range projects {
		urlWithSignature, err := a.publishResourceService(ctx, projects[i].GSThumbnailURL)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		projects[i].SignedThumbnailURL = urlWithSignature
	}

	return projects, nil
}

func (a *customerApplication) GetCountByOpen(ctx context.Context) (int64, error) {
	count, err := a.projectRepository.GetCountByStatus(ctx, pb.Project_Status_Open)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return count, nil
}

func (a *customerApplication) GetCountBySearch(ctx context.Context, params adapter.SearchProjectParams) (int64, error) {
	count, err := a.projectIndexRepository.SearchCountByStatusWithCondition(
		ctx,
		params.Query,
		pb.Project_Status_Open)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return count, nil
}

func (a *customerApplication) GetCountByEntry(ctx context.Context) (int64, error) {
	me, err := a.customerRepository.Get(ctx, a.executorID)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	count, err := a.entryRepository.GetCountByCustomer(ctx, me.ID)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return count, nil
}

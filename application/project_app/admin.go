package project_app

import (
	"context"

	"github.com/pkg/errors"

	"gae-go-recruiting-server/adapter"
	"gae-go-recruiting-server/domain"
)

type adminApplication struct {
	executorID domain.AdminUserID
	*application
}

func (a *adminApplication) GetAllWithPager(ctx context.Context, page int32, offset int32) ([]*domain.Project, error) {
	projects, err := a.projectRepository.GetAllWithPager(ctx, domain.NewPager(page, offset))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return projects, nil
}

func (a *adminApplication) GetAllByFilterWithPager(
	ctx context.Context,
	params adapter.FilterProjectParams,
	page int32,
	offset int32) ([]*domain.Project, error) {
	pager := domain.NewPager(page, offset)

	projectIDsWithHighlights, err := a.projectIndexRepository.SearchByConditionWithPager(
		ctx,
		params.Query,
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
		for _, id := range projectIDsWithHighlights {
			if projectID == id.ID {
				projects[i].Highlights = id.Highlights
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

func (a *adminApplication) GetAllByCompanyWithPager(
	ctx context.Context,
	companyID domain.CompanyID,
	page int32,
	offset int32) ([]*domain.Project, error) {
	projects, err := a.projectRepository.GetAllByCompanyWithPager(ctx, companyID, domain.NewPager(page, offset))
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

func (a *adminApplication) GetTotalCount(ctx context.Context) (int64, error) {
	count, err := a.projectRepository.GetTotalCount(ctx)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return count, nil
}

func (a *adminApplication) GetCountByFilter(ctx context.Context, params adapter.FilterProjectParams) (int64, error) {
	count, err := a.projectIndexRepository.SearchCountByCondition(
		ctx,
		params.Query)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return count, nil
}

func (a *adminApplication) GetCountByCompany(ctx context.Context, companyID domain.CompanyID) (int64, error) {
	count, err := a.projectRepository.GetCountByCompany(ctx, companyID)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return count, nil
}

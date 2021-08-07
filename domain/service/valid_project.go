package service

import (
	"context"

	"gae-go-recruiting-server/adapter"
	"gae-go-recruiting-server/domain"

	"github.com/pkg/errors"
)

func NewValidProjectService(
	companyRepository adapter.CompanyRepository,
	projectRepository adapter.ProjectRepository) adapter.ValidProjectService {
	return func(
		ctx context.Context,
		client *domain.Client,
		projectID domain.ProjectID) (*domain.Project, error) {
		company, err := companyRepository.Get(ctx, client.CompanyID)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		project, err := projectRepository.Get(ctx, projectID)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		if !company.IsIncludeProject(project) {
			return nil, errors.WithStack(domain.ErrInvalidClient)
		}

		return project, nil
	}
}

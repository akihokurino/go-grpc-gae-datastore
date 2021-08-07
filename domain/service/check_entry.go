package service

import (
	"context"

	"gae-go-recruiting-server/adapter"
	"gae-go-recruiting-server/domain"

	"github.com/pkg/errors"
)

func NewCustomerAlreadyEntryToAnyoneService(
	projectRepository adapter.ProjectRepository,
	entryRepository adapter.EntryRepository) adapter.CustomerAlreadyEntryToAnyoneService {
	return func(ctx context.Context, client *domain.Client, customerID domain.CustomerID) (bool, error) {
		projects, err := projectRepository.GetAllByCompany(ctx, client.CompanyID)
		if err != nil {
			return false, errors.WithStack(err)
		}

		entries, err := entryRepository.GetAllByCustomer(ctx, customerID)
		if err != nil {
			return false, errors.WithStack(err)
		}

		return domain.AlreadyEntryFromProjects(entries, projects), nil
	}
}

func NewCustomerAlreadyEntryToThisService(
	entryRepository adapter.EntryRepository) adapter.CustomerAlreadyEntryToThisService {
	return func(
		ctx context.Context,
		project *domain.Project,
		customerID domain.CustomerID) (bool, error) {
		entries, err := entryRepository.GetAllByCustomer(ctx, customerID)
		if err != nil {
			return false, errors.WithStack(err)
		}

		return domain.AlreadyEntryFromProjects(entries, []*domain.Project{project}), nil
	}
}

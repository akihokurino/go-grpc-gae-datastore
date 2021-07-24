package project_app

import (
	"context"
	"sync"
	"time"

	"go.mercari.io/datastore/boom"

	"gae-go-sample/adapter"
	"gae-go-sample/domain"
	pb "gae-go-sample/proto/go/pb"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

type clientApplication struct {
	executorID domain.ClientID
	*application
}

func (a *clientApplication) GetAllWithPager(
	ctx context.Context,
	page int32,
	offset int32) ([]*domain.Project, error) {
	me, err := a.clientRepository.Get(ctx, a.executorID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	company, err := a.companyRepository.Get(ctx, me.CompanyID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	projects, err := a.projectRepository.GetAllByCompanyWithPager(ctx, company.ID, domain.NewPager(page, offset))
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

func (a *clientApplication) GetAllByOpenExcludeAlreadyEntry(
	ctx context.Context,
	customerID domain.CustomerID) ([]*domain.Project, error) {
	me, err := a.clientRepository.Get(ctx, a.executorID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	company, err := a.companyRepository.Get(ctx, me.CompanyID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	projects, err := a.projectRepository.GetAllByCompanyAndStatus(ctx, company.ID, pb.Project_Status_Open)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	noEntries := make([]*domain.Project, 0, len(projects))

	eg := errgroup.Group{}
	mutex := sync.Mutex{}

	for i := range projects {
		project := projects[i]

		eg.Go(func() error {
			already, err := a.customerAlreadyEntryToThisService(ctx, project, customerID)
			if err != nil {
				return err
			}

			if !already {
				mutex.Lock()
				noEntries = append(noEntries, project)
				mutex.Unlock()
			}

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, errors.WithStack(err)
	}

	for i := range noEntries {
		urlWithSignature, err := a.publishResourceService(ctx, noEntries[i].GSThumbnailURL)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		noEntries[i].SignedThumbnailURL = urlWithSignature
	}

	return noEntries, nil
}

func (a *clientApplication) GetCount(ctx context.Context) (int64, error) {
	client, err := a.clientRepository.Get(ctx, a.executorID)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	company, err := a.companyRepository.Get(ctx, client.CompanyID)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	count, err := a.projectRepository.GetCountByCompany(ctx, company.ID)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return count, nil
}

func (a *clientApplication) Create(
	ctx context.Context,
	params adapter.ProjectParams,
	now time.Time) (*domain.Project, error) {
	me, err := a.clientRepository.Get(ctx, a.executorID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	company, err := a.companyRepository.Get(ctx, me.CompanyID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	project, err := domain.NewProject(
		domain.ProjectID(a.idFactory.UUID()),
		company.ID,
		params.Name,
		params.Description,
		params.ThumbnailURL,
		now)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	project.Draft()

	if err := a.transaction(ctx, func(tx *boom.Transaction) error {
		if err := a.projectRepository.Put(tx, project); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, errors.WithStack(err)
	}

	if err := a.projectIndexRepository.Put(ctx, project, company); err != nil {
		return nil, errors.WithStack(err)
	}

	urlWithSignature, err := a.publishResourceService(ctx, project.GSThumbnailURL)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	project.SignedThumbnailURL = urlWithSignature

	return project, nil
}

func (a *clientApplication) Update(
	ctx context.Context,
	id domain.ProjectID,
	params adapter.ProjectParams) (*domain.Project, error) {
	me, err := a.clientRepository.Get(ctx, a.executorID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	project, err := a.validProjectService(ctx, me, id)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err := project.Update(
		params.Name,
		params.Description,
		params.ThumbnailURL); err != nil {
		return nil, errors.WithStack(err)
	}

	company, err := a.companyRepository.Get(ctx, project.CompanyID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err := a.transaction(ctx, func(tx *boom.Transaction) error {
		if err := a.projectRepository.Put(tx, project); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, errors.WithStack(err)
	}

	if err := a.projectIndexRepository.Put(ctx, project, company); err != nil {
		return nil, errors.WithStack(err)
	}

	urlWithSignature, err := a.publishResourceService(ctx, project.GSThumbnailURL)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	project.SignedThumbnailURL = urlWithSignature

	return project, nil
}

func (a *clientApplication) Open(
	ctx context.Context,
	id domain.ProjectID,
	now time.Time) error {
	me, err := a.clientRepository.Get(ctx, a.executorID)
	if err != nil {
		return errors.WithStack(err)
	}

	project, err := a.validProjectService(ctx, me, id)
	if err != nil {
		return errors.WithStack(err)
	}

	project.Open(now)

	company, err := a.companyRepository.Get(ctx, project.CompanyID)
	if err != nil {
		return errors.WithStack(err)
	}

	if err := a.transaction(ctx, func(tx *boom.Transaction) error {
		if err := a.projectRepository.Put(tx, project); err != nil {
			return err
		}

		if err := a.openNoEntrySupportService(
			ctx,
			tx,
			project,
			now); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return errors.WithStack(err)
	}

	if err := a.projectIndexRepository.Put(ctx, project, company); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (a *clientApplication) Draft(ctx context.Context, id domain.ProjectID) error {
	me, err := a.clientRepository.Get(ctx, a.executorID)
	if err != nil {
		return errors.WithStack(err)
	}

	project, err := a.validProjectService(ctx, me, id)
	if err != nil {
		return errors.WithStack(err)
	}

	project.Draft()

	company, err := a.companyRepository.Get(ctx, project.CompanyID)
	if err != nil {
		return errors.WithStack(err)
	}

	if err := a.transaction(ctx, func(tx *boom.Transaction) error {
		if err := a.projectRepository.Put(tx, project); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return errors.WithStack(err)
	}

	if err := a.projectIndexRepository.Put(ctx, project, company); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (a *clientApplication) Close(ctx context.Context, id domain.ProjectID) error {
	me, err := a.clientRepository.Get(ctx, a.executorID)
	if err != nil {
		return errors.WithStack(err)
	}

	project, err := a.validProjectService(ctx, me, id)
	if err != nil {
		return errors.WithStack(err)
	}

	project.Close()

	company, err := a.companyRepository.Get(ctx, project.CompanyID)
	if err != nil {
		return errors.WithStack(err)
	}

	if err := a.transaction(ctx, func(tx *boom.Transaction) error {
		if err := a.projectRepository.Put(tx, project); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return errors.WithStack(err)
	}

	if err := a.projectIndexRepository.Put(ctx, project, company); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (a *clientApplication) Delete(ctx context.Context, id domain.ProjectID) error {
	me, err := a.clientRepository.Get(ctx, a.executorID)
	if err != nil {
		return errors.WithStack(err)
	}

	project, err := a.validProjectService(ctx, me, id)
	if err != nil {
		return errors.WithStack(err)
	}

	if err := a.transaction(ctx, func(tx *boom.Transaction) error {
		if err := a.projectRepository.Delete(tx, project.ID); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return errors.WithStack(err)
	}

	if err := a.projectIndexRepository.Delete(ctx, project.ID); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

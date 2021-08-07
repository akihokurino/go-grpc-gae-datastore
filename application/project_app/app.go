package project_app

import (
	"context"
	"time"

	"go.mercari.io/datastore/boom"

	"gae-go-recruiting-server/adapter"
	"gae-go-recruiting-server/domain"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

type application struct {
	logger                            adapter.CompositeLogger
	projectRepository                 adapter.ProjectRepository
	companyRepository                 adapter.CompanyRepository
	customerRepository                adapter.CustomerRepository
	clientRepository                  adapter.ClientRepository
	projectIndexRepository            adapter.ProjectIndexRepository
	entryRepository                   adapter.EntryRepository
	messageRoomRepository             adapter.MessageRoomRepository
	noEntrySupportRepository          adapter.NoEntrySupportRepository
	noMessageSupportRepository        adapter.NoMessageSupportRepository
	idFactory                         domain.IDFactory
	transaction                       adapter.Transaction
	thresholdProvider                 adapter.ThresholdProvider
	validProjectService               adapter.ValidProjectService
	closeNoEntrySupportService        adapter.CloseNoEntrySupportService
	closeNoMessageSupportService      adapter.CloseNoMessageSupportService
	customerAlreadyEntryToThisService adapter.CustomerAlreadyEntryToThisService
	openNoEntrySupportService         adapter.OpenNoEntrySupportService
	publishResourceService            adapter.PublishResourceService
}

func NewApplication(
	logger adapter.CompositeLogger,
	projectRepository adapter.ProjectRepository,
	companyRepository adapter.CompanyRepository,
	customerRepository adapter.CustomerRepository,
	clientRepository adapter.ClientRepository,
	projectIndexRepository adapter.ProjectIndexRepository,
	entryRepository adapter.EntryRepository,
	messageRoomRepository adapter.MessageRoomRepository,
	noEntrySupportRepository adapter.NoEntrySupportRepository,
	noMessageSupportRepository adapter.NoMessageSupportRepository,
	idFactory domain.IDFactory,
	transaction adapter.Transaction,
	thresholdProvider adapter.ThresholdProvider,
	validProjectService adapter.ValidProjectService,
	closeNoEntrySupportService adapter.CloseNoEntrySupportService,
	closeNoMessageSupportService adapter.CloseNoMessageSupportService,
	customerAlreadyEntryToThisService adapter.CustomerAlreadyEntryToThisService,
	openNoEntrySupportService adapter.OpenNoEntrySupportService,
	publishResourceService adapter.PublishResourceService) adapter.ProjectApplication {
	return &application{
		logger:                            logger,
		projectRepository:                 projectRepository,
		companyRepository:                 companyRepository,
		customerRepository:                customerRepository,
		clientRepository:                  clientRepository,
		projectIndexRepository:            projectIndexRepository,
		entryRepository:                   entryRepository,
		messageRoomRepository:             messageRoomRepository,
		noEntrySupportRepository:          noEntrySupportRepository,
		noMessageSupportRepository:        noMessageSupportRepository,
		idFactory:                         idFactory,
		transaction:                       transaction,
		thresholdProvider:                 thresholdProvider,
		validProjectService:               validProjectService,
		closeNoEntrySupportService:        closeNoEntrySupportService,
		closeNoMessageSupportService:      closeNoMessageSupportService,
		customerAlreadyEntryToThisService: customerAlreadyEntryToThisService,
		openNoEntrySupportService:         openNoEntrySupportService,
		publishResourceService:            publishResourceService,
	}
}

func (a *application) BuildAsAdmin(id domain.AdminUserID) adapter.ProjectApplicationForAdmin {
	return &adminApplication{
		executorID:  id,
		application: a,
	}
}

func (a *application) BuildAsCustomer(id domain.CustomerID) adapter.ProjectApplicationForCustomer {
	return &customerApplication{
		executorID:  id,
		application: a,
	}
}

func (a *application) BuildAsClient(id domain.ClientID) adapter.ProjectApplicationForClient {
	return &clientApplication{
		executorID:  id,
		application: a,
	}
}

func (a *application) GetAllByIDs(ctx context.Context, ids []domain.ProjectID) ([]*domain.Project, error) {
	projects, err := a.projectRepository.GetMulti(ctx, ids)
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

func (a *application) Get(ctx context.Context, id domain.ProjectID) (*domain.Project, error) {
	project, err := a.projectRepository.Get(ctx, id)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	urlWithSignature, err := a.publishResourceService(ctx, project.GSThumbnailURL)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	project.SignedThumbnailURL = urlWithSignature

	return project, nil
}

func (a *application) SupportNoEntry(ctx context.Context, now time.Time) error {
	supports, err := a.noEntrySupportRepository.GetAllByOpened(ctx)
	if err != nil {
		return errors.WithStack(err)
	}

	projectIDs := make([]domain.ProjectID, 0, len(supports))
	for _, support := range supports {
		projectIDs = append(projectIDs, support.ProjectID)
	}

	projects, err := a.projectRepository.GetMulti(ctx, projectIDs)
	if err != nil {
		return errors.WithStack(err)
	}

	eg := errgroup.Group{}

	for i := range projects {
		project := projects[i]

		eg.Go(func() error {
			if domain.IsOvertime(now, project.OpenedAt, a.thresholdProvider.NoEntryDuration) {
				if err := a.transaction(ctx, func(tx *boom.Transaction) error {
					if err := a.closeNoEntrySupportService(
						ctx,
						tx,
						project); err != nil {
						return err
					}

					return nil
				}); err != nil {
					return errors.WithStack(err)
				}
			}

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (a *application) SupportNoMessage(ctx context.Context, now time.Time) error {
	supports, err := a.noMessageSupportRepository.GetAllByOpened(ctx)
	if err != nil {
		return errors.WithStack(err)
	}

	eg := errgroup.Group{}

	for i := range supports {
		support := supports[i]

		eg.Go(func() error {
			roomID := domain.NewMessageRoomID(support.ProjectID, support.CustomerID, support.CompanyID)

			room, err := a.messageRoomRepository.Get(ctx, roomID)
			if err != nil {
				return errors.WithStack(err)
			}

			if domain.IsOvertime(now, room.CreatedAt, a.thresholdProvider.NoMessageDuration) {
				if err := a.transaction(ctx, func(tx *boom.Transaction) error {
					if err := a.closeNoMessageSupportService(
						ctx,
						tx,
						room); err != nil {
						return err
					}

					return nil
				}); err != nil {
					return errors.WithStack(err)
				}
			}

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

package entry_app

import (
	"context"
	"time"

	"go.mercari.io/datastore/boom"

	"gae-go-sample/domain"

	"github.com/pkg/errors"
)

type customerApplication struct {
	executorID domain.CustomerID
	*application
}

func (a *customerApplication) Create(ctx context.Context, projectID domain.ProjectID, now time.Time) (*domain.Entry, error) {
	me, err := a.customerRepository.Get(ctx, a.executorID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	isExists, err := a.entryRepository.Exists(ctx, domain.NewEntryID(a.executorID, projectID))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if isExists {
		return nil, errors.WithStack(domain.ErrEntryAlreadyExists)
	}

	project, err := a.projectRepository.Get(ctx, projectID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	company, err := a.companyRepository.Get(ctx, project.CompanyID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	entry, err := project.Entry(me, now)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	isExists, err = a.messageRoomRepository.Exists(ctx, domain.NewMessageRoomID(project.ID, me.ID, company.ID))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	clients, err := a.clientRepository.GetAllByCompany(ctx, company.ID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if isExists {
		if err := a.transaction(ctx, func(tx *boom.Transaction) error {
			if err := a.entryRepository.Put(tx, entry); err != nil {
				return err
			}

			if err := a.openNoMessageSupportService(
				ctx,
				tx,
				project,
				company,
				me,
				now); err != nil {
				return err
			}

			if err := a.closeNoEntrySupportService(
				ctx,
				tx,
				project); err != nil {
				return err
			}

			if err := a.projectRepository.Put(tx, project); err != nil {
				return err
			}

			return nil
		}); err != nil {
			return nil, errors.WithStack(err)
		}

		return entry, nil
	}

	room, err := me.EnterRoomWith(company, project, now)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err := a.transaction(ctx, func(tx *boom.Transaction) error {
		if err := a.entryRepository.Put(tx, entry); err != nil {
			return err
		}

		if err := a.messageRoomRepository.Put(tx, room); err != nil {
			return err
		}

		if err := a.openNoMessageSupportService(
			ctx,
			tx,
			project,
			company,
			me,
			now); err != nil {
			return err
		}

		if err := a.closeNoEntrySupportService(
			ctx,
			tx,
			project); err != nil {
			return err
		}

		if err := a.projectRepository.Put(tx, project); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, errors.WithStack(err)
	}

	if err := a.rtMessageRoomRepository.Put(ctx, room, me, clients); err != nil {
		return nil, errors.WithStack(err)
	}

	return entry, nil
}

func (a *customerApplication) Delete(ctx context.Context, projectID domain.ProjectID) error {
	if err := a.transaction(ctx, func(tx *boom.Transaction) error {
		if err := a.entryRepository.Delete(tx, domain.NewEntryID(a.executorID, projectID)); err != nil {
			return nil
		}

		return nil
	}); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

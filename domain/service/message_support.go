package service

import (
	"context"
	"time"

	"go.mercari.io/datastore/boom"

	"gae-go-recruiting-server/adapter"
	"gae-go-recruiting-server/domain"

	"github.com/pkg/errors"
)

func NewOpenNoMessageSupportService(noMessageSupportRepository adapter.NoMessageSupportRepository) adapter.OpenNoMessageSupportService {
	return func(
		ctx context.Context,
		tx *boom.Transaction,
		project *domain.Project,
		company *domain.Company,
		customer *domain.Customer,
		now time.Time) error {
		support := project.StartSupportNoMessageWith(company, customer, now)

		if err := noMessageSupportRepository.Put(tx, support); err != nil {
			return err
		}

		return nil
	}
}

func NewCloseNoMessageSupportService(
	noMessageSupportRepository adapter.NoMessageSupportRepository,
	logger adapter.CompositeLogger) adapter.CloseNoMessageSupportService {
	return func(
		ctx context.Context,
		tx *boom.Transaction,
		room *domain.MessageRoom) error {
		support, err := noMessageSupportRepository.Get(ctx, domain.NewNoMessageSupportID(
			room.ProjectID,
			room.CompanyID,
			room.CustomerID))
		if err != nil && !domain.IsNoSuchEntityErr(err) {
			return errors.WithStack(err)
		}
		if err != nil && domain.IsNoSuchEntityErr(err) {
			logger.Error().With(ctx).Printf("no message support is not found %s", room.ID())
			return nil
		}

		support.Close()

		if err := noMessageSupportRepository.Put(tx, support); err != nil {
			return errors.WithStack(err)
		}

		return nil
	}
}

func NewOpenNoEntrySupportService(noEntrySupportRepository adapter.NoEntrySupportRepository) adapter.OpenNoEntrySupportService {
	return func(
		ctx context.Context,
		tx *boom.Transaction,
		project *domain.Project,
		now time.Time) error {
		support := project.StartSupportNoEntry(now)

		if err := noEntrySupportRepository.Put(tx, support); err != nil {
			return err
		}

		return nil
	}
}

func NewCloseNoEntrySupportService(
	noEntrySupportRepository adapter.NoEntrySupportRepository,
	logger adapter.CompositeLogger) adapter.CloseNoEntrySupportService {
	return func(
		ctx context.Context,
		tx *boom.Transaction,
		project *domain.Project) error {
		support, err := noEntrySupportRepository.Get(ctx, project.ID)
		if err != nil && !domain.IsNoSuchEntityErr(err) {
			return errors.WithStack(err)
		}
		if err != nil && domain.IsNoSuchEntityErr(err) {
			logger.Error().With(ctx).Printf("no entry support is not found %s", project.ID)
			return nil
		}

		support.Close()

		if err := noEntrySupportRepository.Put(tx, support); err != nil {
			return errors.WithStack(err)
		}

		return nil
	}
}

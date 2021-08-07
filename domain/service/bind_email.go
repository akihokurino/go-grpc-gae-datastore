package service

import (
	"context"

	"gae-go-recruiting-server/adapter"
	"gae-go-recruiting-server/domain"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

func NewBindCustomerEmailService(
	fireUserRepository adapter.FireUserRepository,
	logger adapter.CompositeLogger) adapter.BindCustomerEmailService {
	return func(ctx context.Context, customers []*domain.Customer) ([]*domain.Customer, error) {
		eg := errgroup.Group{}

		for i := range customers {
			i := i

			eg.Go(func() error {
				fireUser, err := fireUserRepository.Get(ctx, domain.UserID(customers[i].ID))
				if err == nil {
					customers[i].BindEmail(fireUser)
				} else {
					logger.Warn().With(ctx).Printf("failed get firebase user, %s, %#v", customers[i].ID, err)
				}

				return nil
			})
		}

		if err := eg.Wait(); err != nil {
			return nil, errors.WithStack(err)
		}

		return customers, nil
	}
}

func NewBindClientEmailService(
	fireUserRepository adapter.FireUserRepository,
	logger adapter.CompositeLogger) adapter.BindClientEmailService {
	return func(ctx context.Context, clients []*domain.Client) ([]*domain.Client, error) {
		eg := errgroup.Group{}

		for i := range clients {
			i := i

			eg.Go(func() error {
				fireUser, err := fireUserRepository.Get(ctx, domain.UserID(clients[i].ID))
				if err == nil {
					clients[i].BindEmail(fireUser)
				} else {
					logger.Warn().With(ctx).Printf("failed get firebase user, %s, %#v", clients[i].ID, err)
				}

				return nil
			})
		}

		if err := eg.Wait(); err != nil {
			return nil, errors.WithStack(err)
		}

		return clients, nil
	}
}

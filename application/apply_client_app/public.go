package apply_client_app

import (
	"context"
	"time"

	"go.mercari.io/datastore/boom"

	"gae-go-sample/adapter"
	"gae-go-sample/domain"

	"github.com/pkg/errors"
)

type publicApplication struct {
	*application
}

func (a *publicApplication) Create(
	ctx context.Context,
	params adapter.ApplyClientParams,
	now time.Time) (*domain.ApplyClient, error) {
	isExists, err := a.applyClientRepository.Exists(ctx, params.Email)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if isExists {
		return nil, errors.WithStack(domain.ErrEmailAlreadyExists)
	}

	apply, err := domain.NewApplyClient(
		params.Email,
		params.PhoneNumber,
		params.CompanyName,
		params.WebURL,
		params.AccountName,
		params.AccountNameKana,
		params.Position,
		now)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err := a.transaction(ctx, func(tx *boom.Transaction) error {
		if err := a.applyClientRepository.Put(tx, apply); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, errors.WithStack(err)
	}

	return apply, nil
}

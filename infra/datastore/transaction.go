package datastore

import (
	"context"

	"github.com/pkg/errors"

	"go.mercari.io/datastore/boom"

	"gae-go-sample/adapter"
)

func NewTransaction(df adapter.DSFactory) adapter.Transaction {
	return func(ctx context.Context, fn func(tx *boom.Transaction) error) error {
		b := boom.FromClient(ctx, df(ctx))
		if _, err := b.RunInTransaction(fn); err != nil {
			return errors.WithStack(err)
		}
		return nil
	}
}

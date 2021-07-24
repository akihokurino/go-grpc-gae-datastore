package datastore

import (
	"context"
	"reflect"

	"go.mercari.io/datastore/clouddatastore"

	"go.mercari.io/datastore/boom"

	"gae-go-sample/adapter"

	"gae-go-sample/domain"

	"github.com/pkg/errors"

	w "go.mercari.io/datastore"

	"cloud.google.com/go/datastore"
)

func NewDSFactory(projectID string) adapter.DSFactory {
	return func(ctx context.Context) w.Client {
		dc, err := datastore.NewClient(ctx, projectID)
		if err != nil {
			panic(err)
		}

		client, err := clouddatastore.FromClient(ctx, dc)
		if err != nil {
			panic(err)
		}

		return client
	}
}

func NewClient(df adapter.DSFactory) adapter.DSClient {
	return &client{
		df: df,
	}
}

type client struct {
	df adapter.DSFactory
}

func (c *client) GetAll(ctx context.Context, kind string, dst interface{}, orderBy string) error {
	b := boom.FromClient(ctx, c.df(ctx))
	q := b.Client.NewQuery(kind)

	if orderBy != "" {
		q = q.Order(orderBy)
	}

	if _, err := b.GetAll(q, dst); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (c *client) GetByFilter(
	ctx context.Context,
	kind string,
	dst interface{},
	filters map[string]interface{},
	pager *domain.Pager,
	orderBy string) error {
	b := boom.FromClient(ctx, c.df(ctx))
	q := b.Client.NewQuery(kind)

	for k, v := range filters {
		q = q.Filter(k, v)
	}

	if pager != nil {
		q = q.Offset(pager.Offset()).Limit(pager.Limit())
	}

	if orderBy != "" {
		q = q.Order(orderBy)
	}

	if _, err := b.GetAll(q, dst); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (c *client) GetLast(ctx context.Context, kind string, dst interface{}, orderBy string) error {
	b := boom.FromClient(ctx, c.df(ctx))
	q := b.Client.NewQuery(kind).Limit(1)

	if orderBy != "" {
		q = q.Order(orderBy)
	}

	if _, err := b.GetAll(q, dst); err != nil {
		return errors.WithStack(err)
	}

	s := reflect.ValueOf(dereferenceIfPtr(dst))
	switch s.Kind() {
	case reflect.Slice:
		if s.Len() == 0 {
			return errors.WithStack(domain.ErrNoSuchEntity)
		}
	}

	return nil
}

func (c *client) GetLastByFilter(
	ctx context.Context,
	kind string,
	dst interface{},
	filters map[string]interface{},
	orderBy string) error {
	b := boom.FromClient(ctx, c.df(ctx))
	q := b.Client.NewQuery(kind).Limit(1)

	for k, v := range filters {
		q = q.Filter(k, v)
	}

	if orderBy != "" {
		q = q.Order(orderBy)
	}

	if _, err := b.GetAll(q, dst); err != nil {
		return errors.WithStack(err)
	}

	s := reflect.ValueOf(dereferenceIfPtr(dst))
	switch s.Kind() {
	case reflect.Slice:
		if s.Len() == 0 {
			return errors.WithStack(domain.ErrNoSuchEntity)
		}
	}

	return nil
}

func (c *client) Get(ctx context.Context, dst interface{}) error {
	b := boom.FromClient(ctx, c.df(ctx))

	if err := b.Get(dst); err != nil {
		if err == w.ErrNoSuchEntity {
			return errors.WithStack(domain.ErrNoSuchEntity)
		}
		return errors.WithStack(err)
	}

	return nil
}

func (c *client) GetMulti(ctx context.Context, dst interface{}) error {
	b := boom.FromClient(ctx, c.df(ctx))

	if err := b.GetMulti(dst); err != nil {
		multiErr, ok := err.(w.MultiError)
		if !ok {
			return errors.WithStack(err)
		}

		for _, e := range multiErr {
			if e == w.ErrNoSuchEntity {
				return errors.WithStack(domain.ErrNoSuchEntity)
			}
		}

		return errors.WithStack(err)
	}

	return nil
}

func (c *client) GetMultiWithIgnoreError(ctx context.Context, dst interface{}) error {
	b := boom.FromClient(ctx, c.df(ctx))

	if err := b.GetMulti(dst); err != nil {
		_, ok := err.(w.MultiError)
		if !ok {
			return errors.WithStack(err)
		}
	}

	return nil
}

func (c *client) Exists(ctx context.Context, dst interface{}) (bool, error) {
	b := boom.FromClient(ctx, c.df(ctx))

	if err := b.Get(dst); err != nil {
		if err == w.ErrNoSuchEntity {
			return false, nil
		}
		return false, errors.WithStack(err)
	}

	return true, nil
}

func (c *client) GetTotalCount(ctx context.Context, kind string) (int64, error) {
	b := boom.FromClient(ctx, c.df(ctx))
	q := b.Client.NewQuery(kind)

	count, err := b.Count(q)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return int64(count), nil
}

func (c *client) GetCountByFilter(ctx context.Context, kind string, filters map[string]interface{}) (int64, error) {
	b := boom.FromClient(ctx, c.df(ctx))
	q := b.Client.NewQuery(kind)

	for k, v := range filters {
		q = q.Filter(k, v)
	}

	count, err := b.Count(q)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return int64(count), nil
}

func (c *client) Put(tx *boom.Transaction, src interface{}) error {
	if _, err := tx.Put(src); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (c *client) PutMulti(tx *boom.Transaction, src interface{}) error {
	if _, err := tx.PutMulti(src); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (c *client) Delete(tx *boom.Transaction, src interface{}) error {
	if err := tx.Delete(src); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (c *client) DeleteMulti(tx *boom.Transaction, src interface{}) error {
	if err := tx.DeleteMulti(src); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func dereferenceIfPtr(value interface{}) interface{} {
	return reflect.Indirect(reflect.ValueOf(value)).Interface()
}

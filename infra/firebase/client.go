package firebase

import (
	"context"
	"net/url"

	"gae-go-sample/adapter"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"firebase.google.com/go/db"
	"github.com/pkg/errors"
)

func NewClient(rtdbURL *url.URL, adminKey string) adapter.FirebaseClient {
	return &client{
		rtdbURL:  rtdbURL,
		adminKey: adminKey,
	}
}

type client struct {
	rtdbURL  *url.URL
	adminKey string
}

func (c *client) initApp(ctx context.Context) (*firebase.App, error) {
	var app *firebase.App
	var err error

	conf := &firebase.Config{
		DatabaseURL: c.rtdbURL.String(),
	}

	app, err = firebase.NewApp(ctx, conf)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	return app, nil
}

func (c *client) AuthClient(ctx context.Context) (*auth.Client, error) {
	app, err := c.initApp(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	client, err := app.Auth(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return client, nil
}

func (c *client) RTDBClient(ctx context.Context) (*db.Client, error) {
	app, err := c.initApp(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	client, err := app.Database(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return client, nil
}

package fire_user

import (
	"context"

	"gae-go-sample/adapter"
	"gae-go-sample/domain"

	"firebase.google.com/go/auth"
	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/iterator"
)

type repository struct {
	fireClient  adapter.FirebaseClient
	logger      adapter.CompositeLogger
	adminEmails []string
}

func NewRepository(fireClient adapter.FirebaseClient, logger adapter.CompositeLogger, adminEmails []string) adapter.FireUserRepository {
	return &repository{
		fireClient:  fireClient,
		logger:      logger,
		adminEmails: adminEmails,
	}
}

func (r *repository) GetAll(ctx context.Context) ([]*domain.FireUser, error) {
	client, err := r.fireClient.AuthClient(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	users := make([]*domain.FireUser, 0)

	iter := client.Users(ctx, "")
	for {
		userRecord, err := iter.Next()
		if err == iterator.Done {
			return users, nil
		}
		if err != nil {
			return nil, errors.WithStack(err)
		}

		users = append(users, toUserDomainFromExport(userRecord))
	}
}

func (r *repository) Get(ctx context.Context, id domain.UserID) (*domain.FireUser, error) {
	client, err := r.fireClient.AuthClient(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	userRecord, err := client.GetUser(ctx, string(id))
	if err != nil {
		if auth.IsUserNotFound(err) {
			return nil, errors.WithStack(domain.ErrNoSuchEntity)
		}
		return nil, errors.WithStack(err)
	}

	return toUserDomain(userRecord), nil
}

func (r *repository) GetAdmin(ctx context.Context, id domain.AdminUserID) (*domain.FireAdminUser, error) {
	client, err := r.fireClient.AuthClient(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	userRecord, err := client.GetUser(ctx, string(id))
	if err != nil {
		if auth.IsUserNotFound(err) {
			return nil, errors.WithStack(domain.ErrNoSuchEntity)
		}
		return nil, errors.WithStack(err)
	}

	for _, email := range r.adminEmails {
		if email == userRecord.Email {
			return toAdminDomain(userRecord), nil
		}
	}

	return nil, errors.WithStack(domain.ErrNoSuchEntity)
}

func (r *repository) GetByEmail(ctx context.Context, email string) (*domain.FireUser, error) {
	client, err := r.fireClient.AuthClient(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	userRecord, err := client.GetUserByEmail(ctx, email)
	if err != nil {
		if auth.IsUserNotFound(err) {
			return nil, errors.WithStack(domain.ErrNoSuchEntity)
		}
		return nil, errors.WithStack(err)
	}

	return toUserDomain(userRecord), nil
}

func (r *repository) Create(ctx context.Context, email string, password string) (*domain.FireUser, error) {
	client, err := r.fireClient.AuthClient(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	params := (&auth.UserToCreate{}).
		Email(email).
		Password(password).
		EmailVerified(false).
		Disabled(false)

	userRecord, err := client.CreateUser(ctx, params)
	if err != nil {
		r.logger.Error().With(ctx).Printf("failed create user, email = %s", email)

		if auth.IsEmailAlreadyExists(err) {
			return nil, errors.WithStack(domain.ErrEmailAlreadyExists)
		}

		gerr, ok := err.(*googleapi.Error)
		if !ok {
			return nil, errors.WithStack(err)
		}

		if 400 <= gerr.Code && gerr.Code < 500 {
			return nil, errors.WithStack(domain.ErrBadRequest)
		}

		return nil, errors.WithStack(err)
	}

	return toUserDomain(userRecord), nil
}

func (r *repository) Delete(ctx context.Context, id domain.UserID) error {
	client, err := r.fireClient.AuthClient(ctx)
	if err != nil {
		return errors.WithStack(err)
	}

	err = client.DeleteUser(ctx, string(id))
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

package adapter

import (
	"context"

	"gae-go-recruiting-server/domain"

	"firebase.google.com/go/auth"
	"firebase.google.com/go/db"
)

type FirebaseClient interface {
	AuthClient(ctx context.Context) (*auth.Client, error)
	RTDBClient(ctx context.Context) (*db.Client, error)
}

type FireUserRepository interface {
	GetAll(ctx context.Context) ([]*domain.FireUser, error)
	Get(ctx context.Context, id domain.UserID) (*domain.FireUser, error)
	GetAdmin(ctx context.Context, id domain.AdminUserID) (*domain.FireAdminUser, error)
	GetByEmail(ctx context.Context, email string) (*domain.FireUser, error)
	Create(ctx context.Context, email string, password string) (*domain.FireUser, error)
	Delete(ctx context.Context, id domain.UserID) error
}

type RtMessageRoomRepository interface {
	Put(ctx context.Context, room *domain.MessageRoom, customer *domain.Customer, clients []*domain.Client) error
	Delete(ctx context.Context, room *domain.MessageRoom) error
}

type RtMemberRepository interface {
	UpdateCustomer(ctx context.Context, room *domain.MessageRoom, customer *domain.Customer) error
	UpdateClient(ctx context.Context, room *domain.MessageRoom, client *domain.Client) error
}

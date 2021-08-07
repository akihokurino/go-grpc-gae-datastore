package adapter

import (
	"context"
	"net/url"
	"time"

	"go.mercari.io/datastore/boom"

	"gae-go-recruiting-server/domain"
)

type BindCustomerEmailService func(
	ctx context.Context,
	customers []*domain.Customer) ([]*domain.Customer, error)

type BindClientEmailService func(
	ctx context.Context,
	clients []*domain.Client) ([]*domain.Client, error)

type CustomerAlreadyEntryToAnyoneService func(
	ctx context.Context,
	client *domain.Client,
	customerID domain.CustomerID) (bool, error)

type CustomerAlreadyEntryToThisService func(
	ctx context.Context,
	project *domain.Project,
	customerID domain.CustomerID) (bool, error)

type RollbackFireUserService func(ctx context.Context, userID domain.UserID)

type OpenNoMessageSupportService func(
	ctx context.Context,
	tx *boom.Transaction,
	project *domain.Project,
	company *domain.Company,
	customer *domain.Customer,
	now time.Time) error

type CloseNoMessageSupportService func(
	ctx context.Context,
	tx *boom.Transaction,
	room *domain.MessageRoom) error

type OpenNoEntrySupportService func(
	ctx context.Context,
	tx *boom.Transaction,
	project *domain.Project,
	now time.Time) error

type CloseNoEntrySupportService func(
	ctx context.Context,
	tx *boom.Transaction,
	project *domain.Project) error

type ValidCompanyService func(
	ctx context.Context,
	client *domain.Client) (*domain.Company, error)

type ValidProjectService func(
	ctx context.Context,
	client *domain.Client,
	projectID domain.ProjectID) (*domain.Project, error)

type PublishResourceService func(ctx context.Context, resourceURL *url.URL) (*url.URL, error)

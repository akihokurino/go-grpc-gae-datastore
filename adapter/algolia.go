package adapter

import (
	"context"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"

	"gae-go-sample/domain"
	pb "gae-go-sample/proto/go/pb"
)

type AlgoliaClientFactory func(ctx context.Context) *search.Client

type CustomerIndexRepository interface {
	SearchByStatusWithConditionWithPager(
		ctx context.Context,
		query string,
		status pb.Customer_Status,
		pager *domain.Pager,
		orderBy pb.SearchCustomerRequest_OrderBy) ([]domain.CustomerIDWithHighlight, error)
	SearchCountByStatusWithCondition(
		ctx context.Context,
		query string,
		status pb.Customer_Status) (int64, error)
	Put(
		ctx context.Context,
		customer *domain.Customer) error
	Delete(ctx context.Context, customerID domain.CustomerID) error
	DeleteMulti(ctx context.Context, customerIDs []domain.CustomerID) error
	DeleteAll(ctx context.Context) error
}

type ProjectIndexRepository interface {
	SearchByStatusWithConditionWithPager(
		ctx context.Context,
		query string,
		status pb.Project_Status,
		pager *domain.Pager) ([]domain.ProjectIDWithHighlight, error)
	SearchByConditionWithPager(
		ctx context.Context,
		query string,
		pager *domain.Pager) ([]domain.ProjectIDWithHighlight, error)
	SearchCountByStatusWithCondition(
		ctx context.Context,
		query string,
		status pb.Project_Status) (int64, error)
	SearchCountByCondition(
		ctx context.Context,
		query string) (int64, error)
	Put(ctx context.Context, project *domain.Project, company *domain.Company) error
	Delete(ctx context.Context, projectID domain.ProjectID) error
	DeleteMulti(ctx context.Context, projectIDs []domain.ProjectID) error
	DeleteAll(ctx context.Context) error
}

package adapter

import (
	"context"

	"go.mercari.io/datastore/boom"

	pb "gae-go-sample/proto/go/pb"

	"gae-go-sample/domain"

	w "go.mercari.io/datastore"
)

type DSFactory func(ctx context.Context) w.Client

type DSClient interface {
	GetAll(ctx context.Context, kind string, dst interface{}, orderBy string) error
	GetByFilter(
		ctx context.Context,
		kind string,
		dst interface{},
		filters map[string]interface{},
		pager *domain.Pager,
		orderBy string) error
	GetLast(ctx context.Context, kind string, dst interface{}, orderBy string) error
	GetLastByFilter(
		ctx context.Context,
		kind string,
		dst interface{},
		filters map[string]interface{},
		orderBy string) error
	Get(ctx context.Context, dst interface{}) error
	GetMulti(ctx context.Context, dst interface{}) error
	GetMultiWithIgnoreError(ctx context.Context, dst interface{}) error
	Exists(ctx context.Context, dst interface{}) (bool, error)
	GetTotalCount(ctx context.Context, kind string) (int64, error)
	GetCountByFilter(ctx context.Context, kind string, filters map[string]interface{}) (int64, error)
	Put(tx *boom.Transaction, src interface{}) error
	PutMulti(tx *boom.Transaction, src interface{}) error
	Delete(tx *boom.Transaction, src interface{}) error
	DeleteMulti(tx *boom.Transaction, src interface{}) error
}

type Transaction func(ctx context.Context, fn func(tx *boom.Transaction) error) error

type ApplyClientRepository interface {
	GetAllWithPager(ctx context.Context, pager *domain.Pager) ([]*domain.ApplyClient, error)
	GetAllByStatusWithPager(ctx context.Context, status pb.ApplyClient_Status, pager *domain.Pager) ([]*domain.ApplyClient, error)
	Get(ctx context.Context, email domain.ApplyClientID) (*domain.ApplyClient, error)
	GetTotalCount(ctx context.Context) (int64, error)
	GetCountByStatus(ctx context.Context, status pb.ApplyClient_Status) (int64, error)
	Exists(ctx context.Context, email domain.ApplyClientID) (bool, error)
	Put(tx *boom.Transaction, apply *domain.ApplyClient) error
}

type ClientRepository interface {
	GetAll(ctx context.Context) ([]*domain.Client, error)
	GetAllWithPager(ctx context.Context, pager *domain.Pager) ([]*domain.Client, error)
	GetAllByCompany(ctx context.Context, companyID domain.CompanyID) ([]*domain.Client, error)
	GetAllByCompanyWithPager(ctx context.Context, companyID domain.CompanyID, pager *domain.Pager) ([]*domain.Client, error)
	Get(ctx context.Context, id domain.ClientID) (*domain.Client, error)
	GetMulti(ctx context.Context, ids []domain.ClientID) ([]*domain.Client, error)
	GetTotalCount(ctx context.Context) (int64, error)
	GetCountByCompany(ctx context.Context, companyID domain.CompanyID) (int64, error)
	GetCountByCompanyAndRole(ctx context.Context, companyID domain.CompanyID, role pb.Client_Role) (int64, error)
	Put(tx *boom.Transaction, client *domain.Client) error
}

type CompanyRepository interface {
	GetAll(ctx context.Context) ([]*domain.Company, error)
	GetAllWithPager(ctx context.Context, pager *domain.Pager) ([]*domain.Company, error)
	Get(ctx context.Context, id domain.CompanyID) (*domain.Company, error)
	GetMulti(ctx context.Context, ids []domain.CompanyID) ([]*domain.Company, error)
	GetTotalCount(ctx context.Context) (int64, error)
	Put(tx *boom.Transaction, company *domain.Company) error
}

type ContractRepository interface {
	GetAllWithPager(ctx context.Context, pager *domain.Pager) ([]*domain.Contract, error)
	GetAllByCompanyWithPager(ctx context.Context, companyID domain.CompanyID, pager *domain.Pager) ([]*domain.Contract, error)
	Get(ctx context.Context, id domain.ContractID) (*domain.Contract, error)
	GetTotalCount(ctx context.Context) (int64, error)
	GetCountByCompany(ctx context.Context, companyID domain.CompanyID) (int64, error)
	Exists(ctx context.Context, id domain.ContractID) (bool, error)
	Put(tx *boom.Transaction, contract *domain.Contract) error
	Delete(tx *boom.Transaction, id domain.ContractID) error
}

type CustomerRepository interface {
	GetAll(ctx context.Context) ([]*domain.Customer, error)
	GetAllWithPager(ctx context.Context, pager *domain.Pager) ([]*domain.Customer, error)
	Get(ctx context.Context, id domain.CustomerID) (*domain.Customer, error)
	GetMulti(ctx context.Context, ids []domain.CustomerID) ([]*domain.Customer, error)
	GetMultiWithIgnoreNotFound(ctx context.Context, ids []domain.CustomerID) ([]*domain.Customer, error)
	GetTotalCount(ctx context.Context) (int64, error)
	Put(tx *boom.Transaction, customer *domain.Customer) error
}

type EntryRepository interface {
	GetAllWithPager(ctx context.Context, pager *domain.Pager) ([]*domain.Entry, error)
	GetAllByProject(ctx context.Context, projectID domain.ProjectID) ([]*domain.Entry, error)
	GetAllByProjectWithPager(ctx context.Context, projectID domain.ProjectID, pager *domain.Pager) ([]*domain.Entry, error)
	GetAllByCustomerWithPager(ctx context.Context, customerID domain.CustomerID, pager *domain.Pager) ([]*domain.Entry, error)
	GetAllByCustomer(ctx context.Context, customerID domain.CustomerID) ([]*domain.Entry, error)
	Get(ctx context.Context, id domain.EntryID) (*domain.Entry, error)
	GetTotalCount(ctx context.Context) (int64, error)
	GetCountByProject(ctx context.Context, projectID domain.ProjectID) (int64, error)
	GetCountByCustomer(ctx context.Context, customerID domain.CustomerID) (int64, error)
	Exists(ctx context.Context, id domain.EntryID) (bool, error)
	Put(tx *boom.Transaction, entry *domain.Entry) error
	Delete(tx *boom.Transaction, id domain.EntryID) error
}

type MessageRepository interface {
	GetAllByRoom(ctx context.Context, roomID domain.MessageRoomID) ([]*domain.Message, error)
	GetAllByRoomWithPager(ctx context.Context, roomID domain.MessageRoomID, pager *domain.Pager) ([]*domain.Message, error)
	GetLastByRoom(ctx context.Context, roomID domain.MessageRoomID) (*domain.Message, error)
	Get(ctx context.Context, id domain.MessageID) (*domain.Message, error)
	Put(tx *boom.Transaction, message *domain.Message) error
	Delete(tx *boom.Transaction, id domain.MessageID) error
}

type MessageRoomRepository interface {
	GetAll(ctx context.Context) ([]*domain.MessageRoom, error)
	GetAllByCustomer(ctx context.Context, customerID domain.CustomerID) ([]*domain.MessageRoom, error)
	GetAllByCustomerWithPager(ctx context.Context, customerID domain.CustomerID, pager *domain.Pager) ([]*domain.MessageRoom, error)
	GetAllByCompany(ctx context.Context, companyID domain.CompanyID) ([]*domain.MessageRoom, error)
	GetAllByCompanyWithPager(ctx context.Context, companyID domain.CompanyID, pager *domain.Pager) ([]*domain.MessageRoom, error)
	GetLastByProjectAndCustomer(
		ctx context.Context,
		projectID domain.ProjectID,
		customerID domain.CustomerID) (*domain.MessageRoom, error)
	Get(ctx context.Context, id domain.MessageRoomID) (*domain.MessageRoom, error)
	Exists(ctx context.Context, id domain.MessageRoomID) (bool, error)
	Put(tx *boom.Transaction, room *domain.MessageRoom) error
	Delete(tx *boom.Transaction, id domain.MessageRoomID) error
}

type ProjectRepository interface {
	GetAll(ctx context.Context) ([]*domain.Project, error)
	GetAllWithPager(ctx context.Context, pager *domain.Pager) ([]*domain.Project, error)
	GetAllByCompany(ctx context.Context, companyID domain.CompanyID) ([]*domain.Project, error)
	GetAllByCompanyAndStatus(ctx context.Context, companyID domain.CompanyID, status pb.Project_Status) ([]*domain.Project, error)
	GetAllByCompanyWithPager(ctx context.Context, companyID domain.CompanyID, pager *domain.Pager) ([]*domain.Project, error)
	GetAllByStatusWithPager(ctx context.Context, status pb.Project_Status, pager *domain.Pager) ([]*domain.Project, error)
	GetAllByStatusAndEntryNumOrderWithPager(
		ctx context.Context,
		status pb.Project_Status,
		pager *domain.Pager) ([]*domain.Project, error)
	Get(ctx context.Context, id domain.ProjectID) (*domain.Project, error)
	GetMulti(ctx context.Context, ids []domain.ProjectID) ([]*domain.Project, error)
	GetTotalCount(ctx context.Context) (int64, error)
	GetCountByCompany(ctx context.Context, companyID domain.CompanyID) (int64, error)
	GetCountByStatus(ctx context.Context, status pb.Project_Status) (int64, error)
	GetCountByCompanyAndStatus(ctx context.Context, companyID domain.CompanyID, status pb.Project_Status) (int64, error)
	Put(tx *boom.Transaction, project *domain.Project) error
	Delete(tx *boom.Transaction, id domain.ProjectID) error
}

type NoEntrySupportRepository interface {
	GetAllByOpened(ctx context.Context) ([]*domain.NoEntrySupport, error)
	Get(ctx context.Context, projectID domain.ProjectID) (*domain.NoEntrySupport, error)
	Put(tx *boom.Transaction, support *domain.NoEntrySupport) error
}

type NoMessageSupportRepository interface {
	GetAllByOpened(ctx context.Context) ([]*domain.NoMessageSupport, error)
	Get(ctx context.Context, id domain.NoMessageSupportID) (*domain.NoMessageSupport, error)
	Put(tx *boom.Transaction, support *domain.NoMessageSupport) error
}

type UserRepository interface {
	Get(ctx context.Context, id domain.UserID) (*domain.User, error)
	Exists(ctx context.Context, id domain.UserID) (bool, error)
	Put(tx *boom.Transaction, user *domain.User) error
}

package adapter

import (
	"context"
	"net/url"
	"time"

	pb "gae-go-recruiting-server/proto/go/pb"

	"gae-go-recruiting-server/domain"
)

type ApplyClientParams struct {
	Email           domain.ApplyClientID
	PhoneNumber     string
	CompanyName     string
	WebURL          *url.URL
	AccountName     string
	AccountNameKana string
	Position        string
}

type ApplyClientApplication interface {
	BuildAsAdmin(id domain.AdminUserID) ApplyClientApplicationForAdmin
	BuildAsPublic() ApplyClientApplicationForPublic
}

type ApplyClientApplicationForAdmin interface {
	GetAllWithPager(ctx context.Context, page int32, offset int32) ([]*domain.ApplyClient, error)
	GetAllByFilterWithPager(ctx context.Context, status pb.ApplyClient_Status, page int32, offset int32) ([]*domain.ApplyClient, error)
	Get(ctx context.Context, email domain.ApplyClientID) (*domain.ApplyClient, error)
	GetTotalCount(ctx context.Context) (int64, error)
	GetCountByFilter(ctx context.Context, status pb.ApplyClient_Status) (int64, error)
	Accept(
		ctx context.Context,
		email domain.ApplyClientID,
		password string,
		now time.Time) error
	Deny(ctx context.Context, email domain.ApplyClientID) (*domain.ApplyClient, error)
}

type ApplyClientApplicationForPublic interface {
	Create(
		ctx context.Context,
		params ApplyClientParams,
		now time.Time) (*domain.ApplyClient, error)
}

type CreateClientParams struct {
	Name        string
	NameKana    string
	IconURL     *url.URL
	PhoneNumber string
	Position    string
	Role        pb.Client_Role
}

type UpdateClientParams struct {
	Name        string
	NameKana    string
	IconURL     *url.URL
	PhoneNumber string
	Position    string
}

type ClientApplication interface {
	BuildAsAdmin(id domain.AdminUserID) ClientApplicationForAdmin
	BuildAsClient(id domain.ClientID) ClientApplicationForClient
}

type ClientApplicationForAdmin interface {
	GetAllWithPager(ctx context.Context, page int32, offset int32) ([]*domain.Client, error)
	GetAllByCompanyWithPager(ctx context.Context, companyID domain.CompanyID, page int32, offset int32) ([]*domain.Client, error)
	Get(ctx context.Context, id domain.ClientID) (*domain.Client, error)
	GetTotalCount(ctx context.Context) (int64, error)
	GetCountByCompany(ctx context.Context, companyID domain.CompanyID) (int64, error)
}

type ClientApplicationForClient interface {
	GetAll(ctx context.Context, page int32, offset int32) ([]*domain.Client, error)
	GetAllByIDs(ctx context.Context, ids []domain.ClientID) ([]*domain.Client, error)
	Get(ctx context.Context, id domain.ClientID) (*domain.Client, error)
	Create(ctx context.Context, email string, password string, params CreateClientParams, now time.Time) (*domain.Client, error)
	Update(ctx context.Context, params UpdateClientParams) (*domain.Client, error)
	UpdateRole(ctx context.Context, id domain.ClientID, role pb.Client_Role) (*domain.Client, error)
	Delete(ctx context.Context, id domain.ClientID) error
}

type CompanyParams struct {
	Name                       string
	LogoURL                    *url.URL
	WebURL                     *url.URL
	EstablishedAt              time.Time
	PostalCode                 string
	RepresentativeName         string
	CapitalStock               string
	Introduction               string
	AccordingCompanyName       string
	AccordingCompanyPostalCode string
	AccordingCompanyAddress    string
}

type CompanyIDWithProjectID struct {
	CompanyID domain.CompanyID
	ProjectID domain.ProjectID
}

type CompanyWithProjectID struct {
	Company   *domain.Company
	ProjectID domain.ProjectID
}

type CompanyApplication interface {
	BuildAsAdmin(id domain.AdminUserID) CompanyApplicationForAdmin
	BuildAsCustomer(id domain.CustomerID) CompanyApplicationForCustomer
	BuildAsClient(id domain.ClientID) CompanyApplicationForClient
}

type CompanyApplicationForAdmin interface {
	GetAllByIDs(ctx context.Context, ids []domain.CompanyID) ([]*domain.Company, error)
	GetAllWithPager(ctx context.Context, page int32, offset int32) ([]*domain.Company, error)
	Get(ctx context.Context, id domain.CompanyID) (*domain.Company, error)
	GetTotalCount(ctx context.Context) (int64, error)
	Active(ctx context.Context, id domain.CompanyID) error
	Ban(ctx context.Context, id domain.CompanyID) error
}

type CompanyApplicationForCustomer interface {
	GetAllByIDsWithMaskIfNotActive(
		ctx context.Context,
		ids []domain.CompanyID) ([]*domain.Company, error)
	GetAllByIDsWithMaskIfNeedOrNotActive(
		ctx context.Context,
		ids []CompanyIDWithProjectID) ([]*CompanyWithProjectID, error)
	GetWithMaskIfNotActive(ctx context.Context, id domain.CompanyID) (*domain.Company, error)
	GetWithMaskIfNeedOrNotActive(
		ctx context.Context,
		id CompanyIDWithProjectID) (*domain.Company, error)
}

type CompanyApplicationForClient interface {
	GetAllByIDs(ctx context.Context, ids []domain.CompanyID) ([]*domain.Company, error)
	GetAllByIDsWithMaskIfNeed(ctx context.Context, ids []CompanyIDWithProjectID) ([]*CompanyWithProjectID, error)
	Get(ctx context.Context, id domain.CompanyID) (*domain.Company, error)
	GetWithMaskIfNeed(ctx context.Context, id CompanyIDWithProjectID) (*domain.Company, error)
	Update(ctx context.Context, params CompanyParams) (*domain.Company, error)
}

type ContractIDParams struct {
	CompanyID  domain.CompanyID
	ProjectID  domain.ProjectID
	CustomerID domain.CustomerID
}

type ContractIDWithoutCompanyIDParams struct {
	ProjectID  domain.ProjectID
	CustomerID domain.CustomerID
}

type ContractApplication interface {
	BuildAsAdmin(id domain.AdminUserID) ContractApplicationForAdmin
	BuildAsClient(id domain.ClientID) ContractApplicationForClient
}

type ContractApplicationForAdmin interface {
	GetAllWithPager(ctx context.Context, page int32, offset int32) ([]*domain.Contract, error)
	Get(ctx context.Context, idParams ContractIDParams) (*domain.Contract, error)
	GetTotalCount(ctx context.Context) (int64, error)
	Accept(ctx context.Context, idParams ContractIDParams) error
	Cancel(ctx context.Context, idParams ContractIDParams) error
	DownloadFile(ctx context.Context, idParams ContractIDParams) (*domain.File, error)
}

type ContractApplicationForClient interface {
	GetAllWithPager(ctx context.Context, page int32, offset int32) ([]*domain.Contract, error)
	GetAllNewestByIDs(
		ctx context.Context,
		idParams []ContractIDWithoutCompanyIDParams) ([]*domain.Contract, error)
	Get(
		ctx context.Context,
		idParams ContractIDWithoutCompanyIDParams) (*domain.Contract, error)
	GetCount(ctx context.Context) (int64, error)
	Create(
		ctx context.Context,
		idParams ContractIDWithoutCompanyIDParams,
		fileURL *url.URL,
		now time.Time) (*domain.Contract, error)
	Update(
		ctx context.Context,
		idParams ContractIDWithoutCompanyIDParams,
		fileURL *url.URL,
		now time.Time) (*domain.Contract, error)
	Delete(
		ctx context.Context,
		idParams ContractIDWithoutCompanyIDParams) error
}

type CustomerParams struct {
	Name        string
	NameKana    string
	IconURL     *url.URL
	Birthdate   time.Time
	Gender      pb.User_Gender
	PhoneNumber string
	Pr          string
	Address     string
	ResumeURL   *url.URL
}

type SearchCustomerParams struct {
	Query   string
	OrderBy pb.SearchCustomerRequest_OrderBy
}

type FilterCustomerParams struct {
	Query  string
	Status pb.Customer_Status
}

type CustomerApplication interface {
	BuildAsAdmin(id domain.AdminUserID) CustomerApplicationForAdmin
	BuildAsCustomer(id domain.CustomerID) CustomerApplicationForCustomer
	BuildAsClient(id domain.ClientID) CustomerApplicationForClient
}

type CustomerApplicationForAdmin interface {
	GetAll(ctx context.Context) ([]*domain.Customer, error)
	GetAllByIDs(ctx context.Context, ids []domain.CustomerID) ([]*domain.Customer, error)
	GetAllWithPager(ctx context.Context, page int32, offset int32) ([]*domain.Customer, error)
	GetAllByFilterWithPager(
		ctx context.Context,
		params FilterCustomerParams,
		page int32,
		offset int32) ([]*domain.Customer, error)
	GetTotalCount(ctx context.Context) (int64, error)
	GetCountByFilter(ctx context.Context, params FilterCustomerParams) (int64, error)
	Get(ctx context.Context, id domain.CustomerID) (*domain.Customer, error)
	DenyInspection(ctx context.Context, id domain.CustomerID) error
	PassInspection(ctx context.Context, id domain.CustomerID) error
	Deny(ctx context.Context, id domain.CustomerID) error
	ReInspect(ctx context.Context, id domain.CustomerID) error
	Active(ctx context.Context, id domain.CustomerID) error
}

type CustomerApplicationForCustomer interface {
	Create(
		ctx context.Context,
		name string,
		nameKana string,
		gender pb.User_Gender,
		phoneNumber string,
		birthdate time.Time,
		now time.Time) (*domain.Customer, error)
	Update(ctx context.Context, params CustomerParams) (*domain.Customer, error)
}

type CustomerApplicationForClient interface {
	GetAll(ctx context.Context) ([]*domain.Customer, error)
	GetAllByIDs(ctx context.Context, ids []domain.CustomerID) ([]*domain.Customer, error)
	SearchWithPager(
		ctx context.Context,
		params SearchCustomerParams,
		page int32,
		offset int32) ([]*domain.Customer, error)
	GetCountBySearch(ctx context.Context, params SearchCustomerParams) (int64, error)
	Get(ctx context.Context, id domain.CustomerID) (*domain.Customer, error)
}

type EntryApplication interface {
	BuildAsAdmin(id domain.AdminUserID) EntryApplicationForAdmin
	BuildAsCustomer(id domain.CustomerID) EntryApplicationForCustomer
	BuildAsClient(id domain.ClientID) EntryApplicationForClient
}

type EntryApplicationForAdmin interface {
	GetAllWithPager(ctx context.Context, page int32, offset int32) ([]*domain.Entry, error)
	GetAllByProjectWithPager(
		ctx context.Context,
		projectID domain.ProjectID,
		page int32,
		offset int32) ([]*domain.Entry, error)
	Get(ctx context.Context, customerID domain.CustomerID, projectID domain.ProjectID) (*domain.Entry, error)
	GetTotalCount(ctx context.Context) (int64, error)
	GetCountByProject(ctx context.Context, projectID domain.ProjectID) (int64, error)
}

type EntryApplicationForCustomer interface {
	Create(ctx context.Context, projectID domain.ProjectID, now time.Time) (*domain.Entry, error)
	Delete(ctx context.Context, projectID domain.ProjectID) error
}

type EntryApplicationForClient interface {
	GetAllByProject(
		ctx context.Context,
		projectID domain.ProjectID) ([]*domain.Entry, error)
	GetCountByProjects(
		ctx context.Context,
		projectIDs []domain.ProjectID) (map[domain.ProjectID]int64, error)
	GetCountByProject(ctx context.Context, projectID domain.ProjectID) (int64, error)
}

type MessageParams struct {
	Text     string
	ImageURL *url.URL
	FileURL  *url.URL
}

type MessageApplication interface {
	BuildAsAdmin(id domain.AdminUserID) MessageApplicationForAdmin
	BuildAsCustomer(id domain.CustomerID) MessageApplicationForCustomer
	BuildAsClient(id domain.ClientID) MessageApplicationForClient
	Create(
		ctx context.Context,
		messageID domain.MessageID,
		roomID domain.MessageRoomID,
		fromID string,
		toID string,
		params MessageParams,
		now time.Time) error
}

type MessageApplicationForAdmin interface {
	GetAllByRoomWithPager(
		ctx context.Context,
		roomID domain.MessageRoomID,
		page int32,
		offset int32) ([]*domain.Message, error)
	GetAllNewestByRooms(
		ctx context.Context,
		roomIDs []domain.MessageRoomID) ([]*domain.Message, error)
}

type MessageApplicationForCustomer interface {
	GetAllByRoomWithPager(
		ctx context.Context,
		projectID domain.ProjectID,
		companyID domain.CompanyID,
		page int32,
		offset int32) ([]*domain.Message, error)
	GetAllNewestByRooms(
		ctx context.Context,
		roomIDParams []struct {
			ProjectID domain.ProjectID
			CompanyID domain.CompanyID
		}) ([]*domain.Message, error)
}

type MessageApplicationForClient interface {
	GetAllByRoomWithPager(
		ctx context.Context,
		projectID domain.ProjectID,
		customerID domain.CustomerID,
		page int32,
		offset int32) ([]*domain.Message, error)
	GetAllNewestByRooms(
		ctx context.Context,
		roomIDParams []struct {
			ProjectID  domain.ProjectID
			CustomerID domain.CustomerID
		}) ([]*domain.Message, error)
}

type MessageRoomApplication interface {
	BuildAsAdmin(id domain.AdminUserID) MessageRoomApplicationForAdmin
	BuildAsCustomer(id domain.CustomerID) MessageRoomApplicationForCustomer
	BuildAsClient(id domain.ClientID) MessageRoomApplicationForClient
}

type MessageRoomApplicationForAdmin interface {
	GetAllByIDs(ctx context.Context, ids []struct {
		ProjectID  domain.ProjectID
		CustomerID domain.CustomerID
	}) ([]*domain.MessageRoom, error)
}

type MessageRoomApplicationForCustomer interface {
	GetAllWithPager(
		ctx context.Context,
		page int32,
		offset int32) ([]*domain.MessageRoom, error)
	Get(
		ctx context.Context,
		projectID domain.ProjectID,
		companyID domain.CompanyID) (*domain.MessageRoom, error)
	Read(
		ctx context.Context,
		projectID domain.ProjectID,
		companyID domain.CompanyID) error
	Delete(
		ctx context.Context,
		projectID domain.ProjectID,
		companyID domain.CompanyID) error
}

type MessageRoomApplicationForClient interface {
	GetAllWithPager(
		ctx context.Context,
		page int32,
		offset int32) ([]*domain.MessageRoom, error)
	Get(
		ctx context.Context,
		projectID domain.ProjectID,
		customerID domain.CustomerID) (*domain.MessageRoom, error)
	Read(
		ctx context.Context,
		projectID domain.ProjectID,
		customerID domain.CustomerID) error
	Delete(
		ctx context.Context,
		projectID domain.ProjectID,
		customerID domain.CustomerID) error
}

type PrefectureApplication interface {
	GetAll(ctx context.Context) ([]*domain.Prefecture, error)
}

type ProjectParams struct {
	Name         string
	Description  string
	ThumbnailURL *url.URL
}

type SearchProjectParams struct {
	Query string
}

type FilterProjectParams struct {
	Query string
}

type ProjectApplication interface {
	BuildAsAdmin(id domain.AdminUserID) ProjectApplicationForAdmin
	BuildAsCustomer(id domain.CustomerID) ProjectApplicationForCustomer
	BuildAsClient(id domain.ClientID) ProjectApplicationForClient
	SupportNoEntry(ctx context.Context, now time.Time) error
	SupportNoMessage(ctx context.Context, now time.Time) error
}

type ProjectApplicationForAdmin interface {
	GetAllByIDs(ctx context.Context, ids []domain.ProjectID) ([]*domain.Project, error)
	GetAllWithPager(ctx context.Context, page int32, offset int32) ([]*domain.Project, error)
	GetAllByFilterWithPager(ctx context.Context, params FilterProjectParams, page int32, offset int32) ([]*domain.Project, error)
	GetAllByCompanyWithPager(ctx context.Context, companyID domain.CompanyID, page int32, offset int32) ([]*domain.Project, error)
	Get(ctx context.Context, id domain.ProjectID) (*domain.Project, error)
	GetTotalCount(ctx context.Context) (int64, error)
	GetCountByFilter(ctx context.Context, params FilterProjectParams) (int64, error)
	GetCountByCompany(ctx context.Context, companyID domain.CompanyID) (int64, error)
}

type ProjectApplicationForCustomer interface {
	GetAllByIDs(ctx context.Context, ids []domain.ProjectID) ([]*domain.Project, error)
	GetAllByOpenWithPager(ctx context.Context, page int32, offset int32) ([]*domain.Project, error)
	GetAllByRecommend(ctx context.Context) ([]*domain.Project, error)
	SearchWithPager(ctx context.Context, params SearchProjectParams, page int32, offset int32) ([]*domain.Project, error)
	GetAllByEntryWithPager(ctx context.Context, page int32, offset int32) ([]*domain.Project, error)
	Get(ctx context.Context, id domain.ProjectID) (*domain.Project, error)
	GetCountByOpen(ctx context.Context) (int64, error)
	GetCountBySearch(ctx context.Context, params SearchProjectParams) (int64, error)
	GetCountByEntry(ctx context.Context) (int64, error)
	GetAllByEntryNumOrderWithPager(ctx context.Context, page int32, offset int32) ([]*domain.Project, error)
}

type ProjectApplicationForClient interface {
	GetAllByIDs(ctx context.Context, ids []domain.ProjectID) ([]*domain.Project, error)
	GetAllWithPager(ctx context.Context, page int32, offset int32) ([]*domain.Project, error)
	GetAllByOpenExcludeAlreadyEntry(
		ctx context.Context,
		customerID domain.CustomerID) ([]*domain.Project, error)
	Get(ctx context.Context, id domain.ProjectID) (*domain.Project, error)
	GetCount(ctx context.Context) (int64, error)
	Create(ctx context.Context, params ProjectParams, now time.Time) (*domain.Project, error)
	Update(ctx context.Context, id domain.ProjectID, params ProjectParams) (*domain.Project, error)
	Open(ctx context.Context, id domain.ProjectID, now time.Time) error
	Draft(ctx context.Context, id domain.ProjectID) error
	Close(ctx context.Context, id domain.ProjectID) error
	Delete(ctx context.Context, id domain.ProjectID) error
}

type UserApplication interface {
	Get(ctx context.Context, id domain.UserID) (*domain.Me, error)
}

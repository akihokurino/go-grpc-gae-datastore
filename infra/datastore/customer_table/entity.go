package customer_table

import (
	"net/url"
	"time"

	"gae-go-sample/domain"
	pb "gae-go-sample/proto/go/pb"
)

const kind = "Customer"

type entity struct {
	_kind       string `boom:"kind,Customer"`
	ID          string `boom:"id"`
	Name        string
	NameKana    string
	IconURL     string
	Birthdate   time.Time
	Gender      int32
	PhoneNumber string
	Pr          string `datastore:",noindex"`
	Address     string
	Status      int32
	ResumeURL   string
	CreatedAt   time.Time
}

func onlyID(id domain.CustomerID) *entity {
	return &entity{ID: id.String()}
}

func (e *entity) toDomain() *domain.Customer {
	iconURL, _ := url.Parse(e.IconURL)
	resumeURL, _ := url.Parse(e.ResumeURL)

	return &domain.Customer{
		ID:          domain.CustomerID(e.ID),
		Name:        e.Name,
		NameKana:    e.NameKana,
		GSIconURL:   iconURL,
		Birthdate:   e.Birthdate,
		Gender:      pb.User_Gender(e.Gender),
		PhoneNumber: e.PhoneNumber,
		Pr:          e.Pr,
		Address:     e.Address,
		Status:      pb.Customer_Status(e.Status),
		ResumeURL:   resumeURL,
		CreatedAt:   e.CreatedAt,
	}
}

func toEntity(from *domain.Customer) *entity {
	iconURL := ""
	if from.GSIconURL != nil {
		iconURL = from.GSIconURL.String()
	}

	resumeURL := ""
	if from.ResumeURL != nil {
		resumeURL = from.ResumeURL.String()
	}

	return &entity{
		ID:          from.ID.String(),
		Name:        from.Name,
		NameKana:    from.NameKana,
		IconURL:     iconURL,
		Birthdate:   from.Birthdate,
		Gender:      int32(from.Gender),
		PhoneNumber: from.PhoneNumber,
		Pr:          from.Pr,
		Address:     from.Address,
		Status:      int32(from.Status),
		ResumeURL:   resumeURL,
		CreatedAt:   from.CreatedAt,
	}
}

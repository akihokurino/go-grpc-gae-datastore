package company_table

import (
	"net/url"

	"time"

	"gae-go-sample/domain"
	pb "gae-go-sample/proto/go/pb"
)

const kind = "Company"

type entity struct {
	_kind                      string `boom:"kind,Company"`
	ID                         string `boom:"id"`
	Name                       string
	LogoURL                    string
	WebURL                     string
	EstablishedAt              time.Time
	PostalCode                 string
	RepresentativeName         string
	CapitalStock               string
	Introduction               string `datastore:",noindex"`
	AccordingCompanyName       string
	AccordingCompanyPostalCode string
	AccordingCompanyAddress    string
	Status                     int32
	CreatedAt                  time.Time
}

func onlyID(id domain.CompanyID) *entity {
	return &entity{ID: id.String()}
}

func (e *entity) toDomain() *domain.Company {
	logoURL, _ := url.Parse(e.LogoURL)
	webURL, _ := url.Parse(e.WebURL)

	return &domain.Company{
		ID:                         domain.CompanyID(e.ID),
		Name:                       e.Name,
		GSLogoURL:                  logoURL,
		WebURL:                     webURL,
		EstablishedAt:              e.EstablishedAt,
		PostalCode:                 e.PostalCode,
		RepresentativeName:         e.RepresentativeName,
		CapitalStock:               e.CapitalStock,
		Introduction:               e.Introduction,
		AccordingCompanyName:       e.AccordingCompanyName,
		AccordingCompanyPostalCode: e.AccordingCompanyPostalCode,
		AccordingCompanyAddress:    e.AccordingCompanyAddress,
		Status:                     pb.Company_Status(e.Status),
		CreatedAt:                  e.CreatedAt,
	}
}

func toEntity(from *domain.Company) *entity {
	logoURL := ""
	if from.GSLogoURL != nil {
		logoURL = from.GSLogoURL.String()
	}

	webURL := ""
	if from.WebURL != nil {
		webURL = from.WebURL.String()
	}

	return &entity{
		ID:                         from.ID.String(),
		Name:                       from.Name,
		LogoURL:                    logoURL,
		WebURL:                     webURL,
		EstablishedAt:              from.EstablishedAt,
		PostalCode:                 from.PostalCode,
		RepresentativeName:         from.RepresentativeName,
		CapitalStock:               from.CapitalStock,
		Introduction:               from.Introduction,
		AccordingCompanyName:       from.AccordingCompanyName,
		AccordingCompanyPostalCode: from.AccordingCompanyPostalCode,
		AccordingCompanyAddress:    from.AccordingCompanyAddress,
		Status:                     int32(from.Status),
		CreatedAt:                  from.CreatedAt,
	}
}

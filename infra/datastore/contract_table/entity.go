package contract_table

import (
	"net/url"
	"time"

	pb "gae-go-recruiting-server/proto/go/pb"

	"gae-go-recruiting-server/domain"
)

const kind = "Contract"

type entity struct {
	_kind      string `boom:"kind,Contract"`
	ID         string `boom:"id"`
	ProjectID  string
	CompanyID  string
	CustomerID string
	FileURL    string
	Status     int32
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func onlyID(id domain.ContractID) *entity {
	return &entity{ID: id.String()}
}

func (e *entity) toDomain() (*domain.Contract, error) {
	filePath, err := url.Parse(e.FileURL)
	if err != nil {
		return nil, err
	}

	return &domain.Contract{
		ProjectID:  domain.ProjectID(e.ProjectID),
		CompanyID:  domain.CompanyID(e.CompanyID),
		CustomerID: domain.CustomerID(e.CustomerID),
		GSFileURL:  filePath,
		Status:     pb.Contract_Status(e.Status),
		CreatedAt:  e.CreatedAt,
		UpdatedAt:  e.UpdatedAt,
	}, nil
}

func toEntity(from *domain.Contract) *entity {
	return &entity{
		ID:         from.ID().String(),
		ProjectID:  from.ProjectID.String(),
		CompanyID:  from.CompanyID.String(),
		CustomerID: from.CustomerID.String(),
		FileURL:    from.GSFileURL.String(),
		Status:     int32(from.Status),
		CreatedAt:  from.CreatedAt,
		UpdatedAt:  from.UpdatedAt,
	}
}

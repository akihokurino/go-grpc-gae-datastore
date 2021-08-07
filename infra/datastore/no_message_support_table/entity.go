package no_message_support_table

import (
	"time"

	"gae-go-recruiting-server/domain"
)

const kind = "NoMessageSupport"

type entity struct {
	_kind      string `boom:"kind,NoMessageSupport"`
	ID         string `boom:"id"`
	ProjectID  string
	CompanyID  string
	CustomerID string
	Closed     bool
	CreatedAt  time.Time
}

func onlyID(id domain.NoMessageSupportID) *entity {
	return &entity{ID: id.String()}
}

func (e *entity) toDomain() *domain.NoMessageSupport {
	return &domain.NoMessageSupport{
		ProjectID:  domain.ProjectID(e.ProjectID),
		CompanyID:  domain.CompanyID(e.CompanyID),
		CustomerID: domain.CustomerID(e.CustomerID),
		Closed:     e.Closed,
		CreatedAt:  e.CreatedAt,
	}
}

func toEntity(from *domain.NoMessageSupport) *entity {
	return &entity{
		ID:         from.ID().String(),
		ProjectID:  from.ProjectID.String(),
		CompanyID:  from.CompanyID.String(),
		CustomerID: from.CustomerID.String(),
		Closed:     from.Closed,
		CreatedAt:  from.CreatedAt,
	}
}

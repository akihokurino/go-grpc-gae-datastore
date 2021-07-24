package entry_table

import (
	"time"

	"gae-go-sample/domain"
)

const kind = "Entry"

type entity struct {
	_kind      string `boom:"kind,Entry"`
	ID         string `boom:"id"`
	CustomerID string
	ProjectID  string
	CreatedAt  time.Time
}

func onlyID(id domain.EntryID) *entity {
	return &entity{ID: id.String()}
}

func (e *entity) toDomain() *domain.Entry {
	return &domain.Entry{
		CustomerID: domain.CustomerID(e.CustomerID),
		ProjectID:  domain.ProjectID(e.ProjectID),
		CreatedAt:  e.CreatedAt,
	}
}

func toEntity(from *domain.Entry) *entity {
	return &entity{
		ID:         from.ID().String(),
		CustomerID: from.CustomerID.String(),
		ProjectID:  from.ProjectID.String(),
		CreatedAt:  from.CreatedAt,
	}
}

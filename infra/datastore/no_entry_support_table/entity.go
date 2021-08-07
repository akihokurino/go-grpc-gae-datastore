package no_entry_support_table

import (
	"time"

	"gae-go-recruiting-server/domain"
)

const kind = "NoEntrySupport"

type entity struct {
	_kind     string `boom:"kind,NoEntrySupport"`
	ProjectID string `boom:"id"`
	Closed    bool
	CreatedAt time.Time
}

func onlyID(projectID domain.ProjectID) *entity {
	return &entity{ProjectID: projectID.String()}
}

func (e *entity) toDomain() *domain.NoEntrySupport {
	return &domain.NoEntrySupport{
		ProjectID: domain.ProjectID(e.ProjectID),
		Closed:    e.Closed,
		CreatedAt: e.CreatedAt,
	}
}

func toEntity(from *domain.NoEntrySupport) *entity {
	return &entity{
		ProjectID: from.ProjectID.String(),
		Closed:    from.Closed,
		CreatedAt: from.CreatedAt,
	}
}

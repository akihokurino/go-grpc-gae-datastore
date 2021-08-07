package message_room_table

import (
	"time"

	"gae-go-recruiting-server/domain"
)

const kind = "MessageRoom"

type entity struct {
	_kind               string `boom:"kind,MessageRoom"`
	ID                  string `boom:"id"`
	ProjectID           string
	CustomerID          string
	CompanyID           string
	CustomerUnRead      bool
	CustomerUnReadCount int64
	CompanyUnRead       bool
	CompanyUnReadCount  int64
	CreatedAt           time.Time
	CreatedAtYM         string
	UpdatedAt           time.Time
}

func onlyID(id domain.MessageRoomID) *entity {
	return &entity{ID: id.String()}
}

func (e *entity) toDomain() *domain.MessageRoom {
	return &domain.MessageRoom{
		ProjectID:           domain.ProjectID(e.ProjectID),
		CustomerID:          domain.CustomerID(e.CustomerID),
		CompanyID:           domain.CompanyID(e.CompanyID),
		CustomerUnRead:      e.CustomerUnRead,
		CustomerUnReadCount: e.CustomerUnReadCount,
		CompanyUnRead:       e.CompanyUnRead,
		CompanyUnReadCount:  e.CompanyUnReadCount,
		CreatedAt:           e.CreatedAt,
		UpdatedAt:           e.UpdatedAt,
	}
}

func toEntity(from *domain.MessageRoom) *entity {
	createdAtYM := domain.NewYM(from.CreatedAt)

	return &entity{
		ID:                  from.ID().String(),
		ProjectID:           from.ProjectID.String(),
		CustomerID:          from.CustomerID.String(),
		CompanyID:           from.CompanyID.String(),
		CustomerUnRead:      from.CustomerUnRead,
		CustomerUnReadCount: from.CustomerUnReadCount,
		CompanyUnRead:       from.CompanyUnRead,
		CompanyUnReadCount:  from.CompanyUnReadCount,
		CreatedAt:           from.CreatedAt,
		CreatedAtYM:         createdAtYM.String(),
		UpdatedAt:           from.UpdatedAt,
	}
}

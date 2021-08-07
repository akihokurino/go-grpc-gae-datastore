package client_table

import (
	"net/url"
	"time"

	"gae-go-recruiting-server/domain"
	pb "gae-go-recruiting-server/proto/go/pb"
)

const kind = "Client"

type entity struct {
	_kind       string `boom:"kind,Client"`
	ID          string `boom:"id"`
	CompanyID   string
	Name        string
	NameKana    string
	IconURL     string
	PhoneNumber string
	Position    string
	Role        int32
	IsDeleted   bool
	CreatedAt   time.Time
}

func onlyID(id domain.ClientID) *entity {
	return &entity{ID: id.String()}
}

func (e *entity) toDomain() *domain.Client {
	iconURL, _ := url.Parse(e.IconURL)

	return &domain.Client{
		ID:          domain.ClientID(e.ID),
		CompanyID:   domain.CompanyID(e.CompanyID),
		Name:        e.Name,
		NameKana:    e.NameKana,
		GSIconURL:   iconURL,
		PhoneNumber: e.PhoneNumber,
		Position:    e.Position,
		Role:        pb.Client_Role(e.Role),
		IsDeleted:   e.IsDeleted,
		CreatedAt:   e.CreatedAt,
	}
}

func toEntity(from *domain.Client) *entity {
	iconURL := ""
	if from.GSIconURL != nil {
		iconURL = from.GSIconURL.String()
	}

	return &entity{
		ID:          from.ID.String(),
		CompanyID:   from.CompanyID.String(),
		Name:        from.Name,
		NameKana:    from.NameKana,
		IconURL:     iconURL,
		PhoneNumber: from.PhoneNumber,
		Position:    from.Position,
		Role:        int32(from.Role),
		IsDeleted:   from.IsDeleted,
		CreatedAt:   from.CreatedAt,
	}
}

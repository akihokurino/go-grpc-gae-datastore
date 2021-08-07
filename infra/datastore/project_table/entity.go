package project_table

import (
	"net/url"
	"time"

	"gae-go-recruiting-server/domain"
	pb "gae-go-recruiting-server/proto/go/pb"
)

const kind = "Project"

type entity struct {
	_kind        string `boom:"kind,Project"`
	ID           string `boom:"id"`
	CompanyID    string
	Name         string
	Description  string `datastore:",noindex"`
	ThumbnailURL string
	Status       int32
	CreatedAt    time.Time
	OpenedAt     time.Time
}

func onlyID(id domain.ProjectID) *entity {
	return &entity{ID: id.String()}
}

func (e *entity) toDomain() *domain.Project {
	thumbnailURL, _ := url.Parse(e.ThumbnailURL)

	return &domain.Project{
		ID:             domain.ProjectID(e.ID),
		CompanyID:      domain.CompanyID(e.CompanyID),
		Name:           e.Name,
		Description:    e.Description,
		GSThumbnailURL: thumbnailURL,
		Status:         pb.Project_Status(e.Status),
		CreatedAt:      e.CreatedAt,
		OpenedAt:       e.OpenedAt,
	}
}

func toEntity(from *domain.Project) *entity {
	thumbnailURL := ""
	if from.GSThumbnailURL != nil {
		thumbnailURL = from.GSThumbnailURL.String()
	}

	return &entity{
		ID:           from.ID.String(),
		CompanyID:    from.CompanyID.String(),
		Name:         from.Name,
		Description:  from.Description,
		ThumbnailURL: thumbnailURL,
		Status:       int32(from.Status),
		CreatedAt:    from.CreatedAt,
		OpenedAt:     from.OpenedAt,
	}
}

package message_table

import (
	"net/url"
	"time"

	"gae-go-sample/domain"
)

const kind = "Message"

type entity struct {
	_kind       string `boom:"kind,Message"`
	ID          string `boom:"id"`
	RoomID      string
	FromID      string
	ToID        string
	FromCompany bool
	Text        string `datastore:",noindex"`
	ImageURL    string
	FileURL     string
	CreatedAt   time.Time
}

func onlyID(id domain.MessageID) *entity {
	return &entity{ID: id.String()}
}

func (e *entity) toDomain() *domain.Message {
	imageURL, _ := url.Parse(e.ImageURL)
	fileURL, _ := url.Parse(e.FileURL)

	return &domain.Message{
		ID:          domain.MessageID(e.ID),
		RoomID:      domain.MessageRoomID(e.RoomID),
		FromID:      e.FromID,
		ToID:        e.ToID,
		FromCompany: e.FromCompany,
		Text:        e.Text,
		GSImageURL:  imageURL,
		GSFileURL:   fileURL,
		CreatedAt:   e.CreatedAt,
	}
}

func toEntity(from *domain.Message) *entity {
	imageURL := ""
	if from.GSImageURL != nil {
		imageURL = from.GSImageURL.String()
	}

	fileURL := ""
	if from.GSFileURL != nil {
		fileURL = from.GSFileURL.String()
	}

	return &entity{
		ID:          from.ID.String(),
		RoomID:      from.RoomID.String(),
		FromID:      from.FromID,
		ToID:        from.ToID,
		FromCompany: from.FromCompany,
		Text:        from.Text,
		ImageURL:    imageURL,
		FileURL:     fileURL,
		CreatedAt:   from.CreatedAt,
	}
}

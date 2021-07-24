package domain

import (
	"net/url"
	"time"
)

type Message struct {
	ID             MessageID
	RoomID         MessageRoomID
	FromID         string
	ToID           string
	FromCompany    bool
	Text           string
	GSImageURL     *url.URL
	SignedImageURL *url.URL
	GSFileURL      *url.URL
	SignedFileURL  *url.URL
	CreatedAt      time.Time
}

func NewMessage(
	id MessageID,
	roomID MessageRoomID,
	fromID string,
	toID string,
	fromCompany bool,
	text string,
	imageURL *url.URL,
	fileURL *url.URL,
	now time.Time) (*Message, error) {
	if text == "" && imageURL == nil {
		return nil, ErrBadRequest
	}

	return &Message{
		ID:          id,
		RoomID:      roomID,
		FromID:      fromID,
		ToID:        toID,
		FromCompany: fromCompany,
		Text:        text,
		GSImageURL:  imageURL,
		GSFileURL:   fileURL,
		CreatedAt:   now,
	}, nil
}

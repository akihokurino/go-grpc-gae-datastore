package rt_message_room

import (
	"time"
)

type message struct {
	FromUserID string    `json:"fromUserId"`
	ToUserID   string    `json:"toUserId"`
	Text       string    `json:"text"`
	ImageURL   string    `json:"imageUrl"`
	CreatedAt  time.Time `json:"createdAt"`
}

package rt_message_room

import (
	"gae-go-recruiting-server/domain"
)

type messageRoom struct {
	Members  []*member  `json:"members"`
	Messages []*message `json:"messages"`
}

func newMessageRoom() *messageRoom {
	return &messageRoom{
		Members:  make([]*member, 0),
		Messages: make([]*message, 0),
	}
}

func messageRoomPath(roomID domain.MessageRoomID) string {
	return "messageRooms/" + roomID.String()
}

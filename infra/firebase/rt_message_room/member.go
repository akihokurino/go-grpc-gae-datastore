package rt_message_room

import (
	"net/url"

	"gae-go-recruiting-server/domain"
)

type member struct {
	UserID  string `json:"userId"`
	Name    string `json:"name"`
	IconURL string `json:"iconUrl"`
}

func urlString(url *url.URL, alt string) string {
	if url == nil || url.String() == "" {
		return alt
	}
	return url.String()
}

func newMemberFromCustomer(customer *domain.Customer) *member {
	return &member{
		UserID:  string(customer.ID),
		Name:    customer.Name,
		IconURL: urlString(customer.GSIconURL, ""),
	}
}

func newMemberFromClient(client *domain.Client) *member {
	return &member{
		UserID:  string(client.ID),
		Name:    client.Name,
		IconURL: urlString(client.GSIconURL, ""),
	}
}

func memberPath(roomID domain.MessageRoomID, userID domain.UserID) string {
	return messageRoomPath(roomID) + "/members/" + string(userID)
}

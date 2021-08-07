package domain

import (
	pb "gae-go-recruiting-server/proto/go/pb"
)

type User struct {
	ID    UserID
	Email string
	Role  pb.User_Role
}

type FireUser struct {
	UID   UserID
	Email string
}

func FromFireUser(user *FireUser, role pb.User_Role) *User {
	return &User{
		ID:    user.UID,
		Email: user.Email,
		Role:  role,
	}
}

func (u *User) BindEmail(user *FireUser) {
	u.Email = user.Email
}

func (u *User) UpdateRole(role pb.User_Role) {
	u.Role = role
}

func (u *User) IsCustomer() bool {
	return u.Role == pb.User_Role_Customer
}

func (u *User) IsClient() bool {
	return u.Role == pb.User_Role_Client
}

func (u *User) CustomerID() CustomerID {
	if u.IsClient() {
		return ""
	}
	return CustomerID(u.ID)
}

func (u *User) ClientID() ClientID {
	if u.IsCustomer() {
		return ""
	}
	return ClientID(u.ID)
}

type Me struct {
	*User
	Customer         *Customer
	Client           *Client
	MessageStateList []*MessageState
}

type MessageState struct {
	RoomID      MessageRoomID
	IsUnRead    bool
	UnReadCount int64
}

func NewMe(user *User, customer *Customer, client *Client) *Me {
	return &Me{
		User:             user,
		Customer:         customer,
		Client:           client,
		MessageStateList: make([]*MessageState, 0),
	}
}

func (m *Me) DetectUnRead(rooms []*MessageRoom) error {
	stateList := make([]*MessageState, 0, len(rooms))

	var by MessageRoomUser
	switch m.User.Role {
	case pb.User_Role_Customer:
		by = ByCustomer
	case pb.User_Role_Client:
		by = ByCompany
	}

	for _, room := range rooms {
		unRead, err := room.IsUnRead(by)
		if err != nil {
			return err
		}

		unReadCount, err := room.UnReadCount(by)
		if err != nil {
			return err
		}

		stateList = append(stateList, &MessageState{
			RoomID:      room.ID(),
			IsUnRead:    unRead,
			UnReadCount: unReadCount,
		})
	}

	m.MessageStateList = stateList

	return nil
}

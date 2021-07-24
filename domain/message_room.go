package domain

import (
	"time"
)

type MessageRoomUser int

const (
	ByCustomer MessageRoomUser = iota
	ByCompany
)

type MessageRoom struct {
	ProjectID           ProjectID
	CustomerID          CustomerID
	CompanyID           CompanyID
	CustomerUnRead      bool
	CustomerUnReadCount int64
	CompanyUnRead       bool
	CompanyUnReadCount  int64
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

func newMessageRoom(projectID ProjectID, customerID CustomerID, companyID CompanyID, now time.Time) *MessageRoom {
	return &MessageRoom{
		ProjectID:  projectID,
		CustomerID: customerID,
		CompanyID:  companyID,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

func NewMessageRoomID(projectID ProjectID, customerID CustomerID, companyID CompanyID) MessageRoomID {
	return (&MessageRoom{ProjectID: projectID, CustomerID: customerID, CompanyID: companyID}).ID()
}

func (m *MessageRoom) ID() MessageRoomID {
	return MessageRoomID(string(m.ProjectID) + "-" + string(m.CustomerID) + "-" + string(m.CompanyID))
}

func (m *MessageRoom) IsIncludeCustomer(customer *Customer) bool {
	return m.CustomerID == customer.ID
}

func (m *MessageRoom) IsIncludeClient(client *Client) bool {
	return m.CompanyID == client.CompanyID
}

func (m *MessageRoom) IsUnRead(by MessageRoomUser) (bool, error) {
	switch by {
	case ByCustomer:
		return m.CustomerUnRead, nil
	case ByCompany:
		return m.CompanyUnRead, nil
	default:
		return false, ErrInvalidMessageRoomUser
	}
}

func (m *MessageRoom) UnReadCount(by MessageRoomUser) (int64, error) {
	switch by {
	case ByCustomer:
		return m.CustomerUnReadCount, nil
	case ByCompany:
		return m.CompanyUnReadCount, nil
	default:
		return 0, ErrInvalidMessageRoomUser
	}
}

func (m *MessageRoom) ReceiveMessage(by MessageRoomUser, now time.Time) error {
	switch by {
	case ByCustomer:
		m.CustomerUnRead = true
		m.CustomerUnReadCount = m.CustomerUnReadCount + 1
	case ByCompany:
		m.CompanyUnRead = true
		m.CompanyUnReadCount = m.CompanyUnReadCount + 1
	default:
		return ErrInvalidMessageRoomUser
	}

	m.UpdatedAt = now

	return nil
}

func (m *MessageRoom) Read(by MessageRoomUser) error {
	switch by {
	case ByCustomer:
		m.CustomerUnRead = false
		m.CustomerUnReadCount = 0
	case ByCompany:
		m.CompanyUnRead = false
		m.CompanyUnReadCount = 0
	default:
		return ErrInvalidMessageRoomUser
	}

	return nil
}

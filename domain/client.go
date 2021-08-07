package domain

import (
	"net/url"
	"time"

	"gae-go-recruiting-server/util/validator"

	pb "gae-go-recruiting-server/proto/go/pb"
)

type Client struct {
	ID            ClientID
	CompanyID     CompanyID
	Name          string
	NameKana      string
	GSIconURL     *url.URL
	SignedIconURL *url.URL
	PhoneNumber   string
	Position      string
	Email         string
	Role          pb.Client_Role
	IsDeleted     bool
	CreatedAt     time.Time
}

func NewClient(
	user *User,
	companyID CompanyID,
	name string,
	nameKana string,
	phoneNumber string,
	position string,
	role pb.Client_Role,
	now time.Time) (*Client, error) {
	if user.Role != pb.User_Role_Client {
		return nil, ErrInvalidUserRole
	}
	if err := validator.ValidateTextRange(name, 1, 100); err != nil {
		return nil, NewBadRequestError(err.Error())
	}
	if err := validator.ValidateTextRange(nameKana, 1, 120); err != nil {
		return nil, NewBadRequestError(err.Error())
	}
	if err := validator.ValidatePhoneNumber(phoneNumber); err != nil {
		return nil, NewBadRequestError(err.Error())
	}

	return &Client{
		ID:          ClientID(user.ID),
		CompanyID:   companyID,
		Name:        name,
		NameKana:    nameKana,
		PhoneNumber: phoneNumber,
		Position:    position,
		Email:       user.Email,
		Role:        role,
		CreatedAt:   now,
	}, nil
}

func (c *Client) Update(
	name string,
	nameKana string,
	iconURL *url.URL,
	phoneNumber string,
	position string) error {
	if err := validator.ValidateTextRange(name, 1, 100); err != nil {
		return NewBadRequestError(err.Error())
	}
	if err := validator.ValidateTextRange(nameKana, 1, 120); err != nil {
		return NewBadRequestError(err.Error())
	}
	if err := validator.ValidatePhoneNumber(phoneNumber); err != nil {
		return NewBadRequestError(err.Error())
	}
	if err := validator.ValidateTextRange(position, 0, 100); err != nil {
		return NewBadRequestError(err.Error())
	}

	c.Name = name
	c.NameKana = nameKana
	c.GSIconURL = iconURL
	c.PhoneNumber = phoneNumber
	c.Position = position

	return nil
}

func (c *Client) UpdateRole(role pb.Client_Role) {
	c.Role = role
}

func (c *Client) BindEmail(user *FireUser) {
	c.Email = user.Email
}

func (c *Client) IsAdmin() bool {
	return c.Role == pb.Client_Role_Admin
}

func (c *Client) Delete() {
	c.IsDeleted = true
}

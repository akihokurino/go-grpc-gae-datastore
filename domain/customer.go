package domain

import (
	"net/url"
	"time"

	"gae-go-sample/util/validator"

	pb "gae-go-sample/proto/go/pb"
)

type CustomerIDWithHighlight struct {
	ID         CustomerID
	Highlights []*SearchHighlight
}

type Customer struct {
	ID            CustomerID
	Name          string
	NameKana      string
	GSIconURL     *url.URL
	SignedIconURL *url.URL
	Birthdate     time.Time
	Gender        pb.User_Gender
	PhoneNumber   string
	Email         string
	Pr            string
	Address       string
	Status        pb.Customer_Status
	ResumeURL     *url.URL
	CreatedAt     time.Time

	Highlights []*SearchHighlight
}

func NewDefaultCustomer(
	user *User,
	name string,
	nameKana string,
	gender pb.User_Gender,
	phoneNumber string,
	birthdate time.Time,
	now time.Time) (*Customer, error) {
	if user.Role != pb.User_Role_Customer {
		return nil, ErrInvalidUserRole
	}
	if err := validator.ValidateTextRange(name, 1, 100); err != nil {
		return nil, NewBadRequestError(err.Error())
	}
	if err := validator.ValidateTextRange(nameKana, 1, 100); err != nil {
		return nil, NewBadRequestError(err.Error())
	}
	if err := validator.ValidateHiragana(nameKana); err != nil {
		return nil, NewBadRequestError(err.Error())
	}
	if err := validator.ValidatePhoneNumber(phoneNumber); err != nil {
		return nil, NewBadRequestError(err.Error())
	}

	return &Customer{
		ID:          CustomerID(user.ID),
		Name:        name,
		NameKana:    nameKana,
		Birthdate:   birthdate,
		Gender:      gender,
		PhoneNumber: phoneNumber,
		Email:       user.Email,
		Status:      pb.Customer_Status_Inspection,
		CreatedAt:   now,
	}, nil
}

func (c *Customer) Update(
	name string,
	nameKana string,
	iconURL *url.URL,
	birthdate time.Time,
	gender pb.User_Gender,
	phoneNumber string,
	pr string,
	address string,
	resumeURL *url.URL) error {
	if err := validator.ValidateTextRange(name, 1, 100); err != nil {
		return NewBadRequestError(err.Error())
	}
	if err := validator.ValidateTextRange(nameKana, 1, 100); err != nil {
		return NewBadRequestError(err.Error())
	}
	if err := validator.ValidateHiragana(nameKana); err != nil {
		return NewBadRequestError(err.Error())
	}
	if err := validator.ValidatePhoneNumber(phoneNumber); err != nil {
		return NewBadRequestError(err.Error())
	}
	if err := validator.ValidateTextRange(pr, 0, 5000); err != nil {
		return NewBadRequestError(err.Error())
	}

	c.Name = name
	c.NameKana = nameKana
	c.GSIconURL = iconURL
	c.Birthdate = birthdate
	c.Gender = gender
	c.PhoneNumber = phoneNumber
	c.Pr = pr
	c.Address = address
	c.ResumeURL = resumeURL

	return nil
}

func (c *Customer) ReInspect() {
	c.Status = pb.Customer_Status_Inspection
}

func (c *Customer) Deny() error {
	if !c.IsInspection() {
		return ErrBadRequest
	}
	c.Status = pb.Customer_Status_Denied
	return nil
}

func (c *Customer) Active() error {
	if !c.IsInspection() {
		return ErrBadRequest
	}
	c.Status = pb.Customer_Status_Active
	return nil
}

func (c *Customer) IsInspection() bool {
	return c.Status == pb.Customer_Status_Inspection
}

func (c *Customer) IsDenied() bool {
	return c.Status == pb.Customer_Status_Denied
}

func (c *Customer) IsActive() bool {
	return c.Status == pb.Customer_Status_Active
}

func (c *Customer) BindEmail(user *FireUser) {
	c.Email = user.Email
}

func (c *Customer) EnterRoomWith(company *Company, project *Project, now time.Time) (*MessageRoom, error) {
	if c.IsDenied() {
		return nil, ErrCustomerIsNotActive
	}

	if !project.IsOpen() {
		return nil, ErrProjectIsNotOpen
	}

	return newMessageRoom(project.ID, c.ID, company.ID, now), nil
}

func AlreadyEntryFromProjects(entries []*Entry, projects []*Project) bool {
	already := false

LABEL:
	for _, entry := range entries {
		for _, project := range projects {
			if entry.ProjectID == project.ID {
				already = true
				break LABEL
			}
		}
	}

	return already
}

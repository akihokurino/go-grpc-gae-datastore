package domain

import (
	"net/url"
	"regexp"
	"time"

	"gae-go-recruiting-server/util/validator"

	pb "gae-go-recruiting-server/proto/go/pb"
)

type Company struct {
	ID                         CompanyID
	Name                       string
	GSLogoURL                  *url.URL
	SignedLogoURL              *url.URL
	WebURL                     *url.URL
	EstablishedAt              time.Time
	PostalCode                 string
	RepresentativeName         string
	CapitalStock               string
	Introduction               string
	AccordingCompanyName       string
	AccordingCompanyPostalCode string
	AccordingCompanyAddress    string
	Status                     pb.Company_Status
	CreatedAt                  time.Time
}

func newDefaultCompany(id CompanyID, name string, webURL *url.URL, now time.Time) (*Company, error) {
	if err := validator.ValidateTextRange(name, 1, 200); err != nil {
		return nil, NewBadRequestError(err.Error())
	}

	return &Company{
		ID:        id,
		Name:      name,
		WebURL:    webURL,
		Status:    pb.Company_Status_Active,
		CreatedAt: now,
	}, nil
}

func (c *Company) Update(
	name string,
	logoURL *url.URL,
	webURL *url.URL,
	establishedAt time.Time,
	postalCode string,
	representativeName string,
	capitalStock string,
	introduction string,
	accordingCompanyName string,
	accordingCompanyPostalCode string,
	accordingCompanyAddress string) error {
	if err := validator.ValidateTextRange(name, 1, 200); err != nil {
		return NewBadRequestError(err.Error())
	}
	if postalCode != "" {
		if err := validator.ValidatePostalCode(postalCode); err != nil {
			return NewBadRequestError(err.Error())
		}
	}
	if err := validator.ValidateTextRange(representativeName, 0, 100); err != nil {
		return NewBadRequestError(err.Error())
	}
	if capitalStock != "" {
		rep := regexp.MustCompile(`^[1-9]\d*$`)
		if !rep.MatchString(capitalStock) {
			return ErrBadRequest
		}
	}
	if err := validator.ValidateTextRange(introduction, 0, 5000); err != nil {
		return NewBadRequestError(err.Error())
	}
	if err := validator.ValidateTextRange(accordingCompanyName, 0, 200); err != nil {
		return NewBadRequestError(err.Error())
	}
	if accordingCompanyPostalCode != "" {
		if err := validator.ValidatePostalCode(accordingCompanyPostalCode); err != nil {
			return NewBadRequestError(err.Error())
		}
	}
	if err := validator.ValidateTextRange(accordingCompanyAddress, 0, 1000); err != nil {
		return NewBadRequestError(err.Error())
	}

	c.Name = name
	c.GSLogoURL = logoURL
	c.WebURL = webURL
	c.EstablishedAt = establishedAt
	c.PostalCode = postalCode
	c.RepresentativeName = representativeName
	c.CapitalStock = capitalStock
	c.Introduction = introduction
	c.AccordingCompanyName = accordingCompanyName
	c.AccordingCompanyPostalCode = accordingCompanyPostalCode
	c.AccordingCompanyAddress = accordingCompanyAddress

	return nil
}

func (c *Company) Active() {
	c.Status = pb.Company_Status_Active
}

func (c *Company) Ban() {
	c.Status = pb.Company_Status_BAN
}

func (c *Company) IsBan() bool {
	return c.Status == pb.Company_Status_BAN
}

func (c *Company) IsIncludeProject(project *Project) bool {
	return c.ID == project.CompanyID
}

func (c *Company) IsIncludeClient(client *Client) bool {
	return c.ID == client.CompanyID
}

func (c *Company) Contract(project *Project, customer *Customer, fileURL *url.URL, now time.Time) (*Contract, error) {
	if !project.IsOpen() {
		return nil, ErrProjectIsNotOpen
	}

	if !customer.IsActive() {
		return nil, ErrCustomerIsNotActive
	}

	return newContract(project.ID, c.ID, customer.ID, fileURL, now), nil
}

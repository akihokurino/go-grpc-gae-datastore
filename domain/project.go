package domain

import (
	"net/url"
	"time"

	pb "gae-go-sample/proto/go/pb"
	"gae-go-sample/util/validator"
)

type ProjectIDWithHighlight struct {
	ID         ProjectID
	Highlights []*SearchHighlight
}

type Project struct {
	ID                 ProjectID
	CompanyID          CompanyID
	Name               string
	Description        string
	GSThumbnailURL     *url.URL
	SignedThumbnailURL *url.URL
	Status             pb.Project_Status
	CreatedAt          time.Time
	OpenedAt           time.Time

	Highlights []*SearchHighlight
}

func NewProject(
	id ProjectID,
	companyID CompanyID,
	name string,
	description string,
	thumbnailURL *url.URL,
	now time.Time) (*Project, error) {
	if err := validator.ValidateTextRange(name, 1, 40); err != nil {
		return nil, err
	}
	if err := validator.ValidateTextRange(description, 0, 5000); err != nil {
		return nil, err
	}

	return &Project{
		ID:             id,
		CompanyID:      companyID,
		Name:           name,
		Description:    description,
		GSThumbnailURL: thumbnailURL,
		Status:         pb.Project_Status_Open,
		CreatedAt:      now,
	}, nil
}

func (p *Project) Update(
	name string,
	description string,
	thumbnailURL *url.URL) error {
	if err := validator.ValidateTextRange(name, 1, 40); err != nil {
		return err
	}
	if err := validator.ValidateTextRange(description, 0, 5000); err != nil {
		return err
	}

	p.Name = name
	p.Description = description
	p.GSThumbnailURL = thumbnailURL

	return nil
}

func (p *Project) Open(now time.Time) {
	p.Status = pb.Project_Status_Open
	p.OpenedAt = now
}

func (p *Project) Draft() {
	p.Status = pb.Project_Status_Draft
}

func (p *Project) Close() {
	p.Status = pb.Project_Status_Close
}

func (p *Project) IsOpen() bool {
	return p.Status == pb.Project_Status_Open
}

func (p *Project) IsDraft() bool {
	return p.Status == pb.Project_Status_Draft
}

func (p *Project) IsClose() bool {
	return p.Status == pb.Project_Status_Close
}

func (p *Project) Entry(customer *Customer, now time.Time) (*Entry, error) {
	if !customer.IsActive() {
		return nil, ErrCustomerIsNotActive
	}

	if !p.IsOpen() {
		return nil, ErrProjectIsNotOpen
	}

	return newEntry(customer.ID, p.ID, now), nil
}

func (p *Project) StartSupportNoEntry(now time.Time) *NoEntrySupport {
	return &NoEntrySupport{
		ProjectID: p.ID,
		Closed:    false,
		CreatedAt: now,
	}
}

func (p *Project) StartSupportNoMessageWith(
	company *Company,
	customer *Customer,
	now time.Time) *NoMessageSupport {
	return &NoMessageSupport{
		ProjectID:  p.ID,
		CompanyID:  company.ID,
		CustomerID: customer.ID,
		Closed:     false,
		CreatedAt:  now,
	}
}

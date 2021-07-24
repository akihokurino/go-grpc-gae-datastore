package domain

import (
	"fmt"
	"time"
)

type NoEntrySupport struct {
	ProjectID ProjectID
	Closed    bool
	CreatedAt time.Time
}

func (s *NoEntrySupport) Close() {
	s.Closed = true
}

type NoMessageSupport struct {
	ProjectID  ProjectID
	CompanyID  CompanyID
	CustomerID CustomerID
	Closed     bool
	CreatedAt  time.Time
}

func NewNoMessageSupportID(projectID ProjectID, companyID CompanyID, customerID CustomerID) NoMessageSupportID {
	return (&NoMessageSupport{ProjectID: projectID, CompanyID: companyID, CustomerID: customerID}).ID()
}

func (s *NoMessageSupport) ID() NoMessageSupportID {
	return NoMessageSupportID(fmt.Sprintf("%s-%s-%s", string(s.ProjectID), string(s.CompanyID), string(s.CustomerID)))
}

func (s *NoMessageSupport) Close() {
	s.Closed = true
}

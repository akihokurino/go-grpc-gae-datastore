package domain

import (
	"net/url"
	"time"

	pb "gae-go-recruiting-server/proto/go/pb"
)

type Contract struct {
	ProjectID     ProjectID
	CompanyID     CompanyID
	CustomerID    CustomerID
	GSFileURL     *url.URL
	SignedFileURL *url.URL
	Status        pb.Contract_Status
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func newContract(
	projectID ProjectID,
	companyID CompanyID,
	customerID CustomerID,
	fileURL *url.URL,
	now time.Time) *Contract {
	return &Contract{
		ProjectID:  projectID,
		CompanyID:  companyID,
		CustomerID: customerID,
		GSFileURL:  fileURL,
		Status:     pb.Contract_Status_InProgress,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

func NewContractID(projectID ProjectID, companyID CompanyID, customerID CustomerID) ContractID {
	return (&Contract{ProjectID: projectID, CompanyID: companyID, CustomerID: customerID}).ID()
}

func (c *Contract) ID() ContractID {
	return ContractID(string(c.ProjectID) + "-" + string(c.CompanyID) + "-" + string(c.CustomerID))
}

func (c *Contract) Update(fileURL *url.URL, now time.Time) error {
	if c.IsAccepted() {
		return ErrContractAlreadyAccepted
	}

	if c.IsCanceled() {
		return ErrContractAlreadyCanceled
	}

	c.GSFileURL = fileURL
	c.UpdatedAt = now

	return nil
}

func (c *Contract) Accept() error {
	if c.IsCanceled() {
		return ErrContractAlreadyCanceled
	}

	c.Status = pb.Contract_Status_Accepted

	return nil
}

func (c *Contract) Cancel() error {
	if c.IsInProgress() {
		return ErrContractInProgress
	}

	c.Status = pb.Contract_Status_Canceled

	return nil
}

func (c *Contract) IsInProgress() bool {
	return c.Status == pb.Contract_Status_InProgress
}

func (c *Contract) IsAccepted() bool {
	return c.Status == pb.Contract_Status_Accepted
}

func (c *Contract) IsCanceled() bool {
	return c.Status == pb.Contract_Status_Canceled
}

package domain

import (
	"net/url"
	"time"

	pb "gae-go-sample/proto/go/pb"
	"gae-go-sample/util/validator"
)

type ApplyClient struct {
	Email           ApplyClientID
	PhoneNumber     string
	CompanyName     string
	WebURL          *url.URL
	AccountName     string
	AccountNameKana string
	Position        string
	Status          pb.ApplyClient_Status
	CreatedAt       time.Time
}

func NewApplyClient(
	email ApplyClientID,
	phoneNumber string,
	companyName string,
	webURL *url.URL,
	accountName string,
	accountNameKana string,
	position string,
	now time.Time) (*ApplyClient, error) {
	if err := validator.ValidateEmail(email.String()); err != nil {
		return nil, NewBadRequestError(err.Error())
	}
	if err := validator.ValidatePhoneNumber(phoneNumber); err != nil {
		return nil, NewBadRequestError(err.Error())
	}
	if err := validator.ValidateTextRange(companyName, 1, 200); err != nil {
		return nil, NewBadRequestError(err.Error())
	}
	if err := validator.ValidateTextRange(accountName, 1, 100); err != nil {
		return nil, NewBadRequestError(err.Error())
	}
	if err := validator.ValidateTextRange(accountNameKana, 1, 120); err != nil {
		return nil, NewBadRequestError(err.Error())
	}

	return &ApplyClient{
		Email:           email,
		PhoneNumber:     phoneNumber,
		CompanyName:     companyName,
		WebURL:          webURL,
		AccountName:     accountName,
		AccountNameKana: accountNameKana,
		Position:        position,
		Status:          pb.ApplyClient_Status_Inspection,
		CreatedAt:       now,
	}, nil
}

func (a *ApplyClient) Accept() {
	a.Status = pb.ApplyClient_Status_Accepted
}

func (a *ApplyClient) Deny() {
	a.Status = pb.ApplyClient_Status_Denied
}

func (a *ApplyClient) IsAccepted() bool {
	return a.Status == pb.ApplyClient_Status_Accepted
}

func (a *ApplyClient) IsDenied() bool {
	return a.Status == pb.ApplyClient_Status_Denied
}

func (a *ApplyClient) CreateCompanyWithClient(companyID CompanyID, user *User, now time.Time) (*Company, *Client, error) {
	if !a.IsAccepted() {
		return nil, nil, ErrApplyClientNotAccepted
	}

	company, err := newDefaultCompany(companyID, a.CompanyName, a.WebURL, now)
	if err != nil {
		return nil, nil, err
	}

	client, err := NewClient(
		user,
		companyID,
		a.AccountName,
		a.AccountNameKana,
		a.PhoneNumber,
		a.Position,
		pb.Client_Role_Admin,
		now)
	if err != nil {
		return nil, nil, err
	}

	return company, client, nil
}

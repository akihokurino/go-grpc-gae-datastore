package apply_client_table

import (
	"net/url"
	"time"

	"gae-go-recruiting-server/domain"
	pb "gae-go-recruiting-server/proto/go/pb"
)

const kind = "ApplyClient"

type entity struct {
	_kind           string `boom:"kind,ApplyClient"`
	Email           string `boom:"id"`
	PhoneNumber     string
	CompanyName     string
	WebURL          string
	AccountName     string
	AccountNameKana string
	Position        string
	Status          int32
	CreatedAt       time.Time
}

func onlyID(email domain.ApplyClientID) *entity {
	return &entity{Email: email.String()}
}

func (e *entity) toDomain() *domain.ApplyClient {
	webURL, _ := url.Parse(e.WebURL)

	return &domain.ApplyClient{
		Email:           domain.ApplyClientID(e.Email),
		PhoneNumber:     e.PhoneNumber,
		CompanyName:     e.CompanyName,
		WebURL:          webURL,
		AccountName:     e.AccountName,
		AccountNameKana: e.AccountNameKana,
		Position:        e.Position,
		Status:          pb.ApplyClient_Status(e.Status),
		CreatedAt:       e.CreatedAt,
	}
}

func toEntity(from *domain.ApplyClient) *entity {
	return &entity{
		Email:           from.Email.String(),
		PhoneNumber:     from.PhoneNumber,
		CompanyName:     from.CompanyName,
		WebURL:          from.WebURL.String(),
		AccountName:     from.AccountName,
		AccountNameKana: from.AccountNameKana,
		Position:        from.Position,
		Status:          int32(from.Status),
		CreatedAt:       from.CreatedAt,
	}
}

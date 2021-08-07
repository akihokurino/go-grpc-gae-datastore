package fire_user

import (
	"gae-go-recruiting-server/domain"

	"firebase.google.com/go/auth"
)

func toUserDomain(record *auth.UserRecord) *domain.FireUser {
	return &domain.FireUser{
		UID:   domain.UserID(record.UID),
		Email: record.Email,
	}
}

func toUserDomainFromExport(record *auth.ExportedUserRecord) *domain.FireUser {
	return &domain.FireUser{
		UID:   domain.UserID(record.UID),
		Email: record.Email,
	}
}

func toAdminDomain(record *auth.UserRecord) *domain.FireAdminUser {
	return &domain.FireAdminUser{
		UID:   domain.AdminUserID(record.UID),
		Email: record.Email,
	}
}

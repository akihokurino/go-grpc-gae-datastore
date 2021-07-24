package user_table

import (
	"gae-go-sample/domain"
	pb "gae-go-sample/proto/go/pb"
)

const kind = "User"

type entity struct {
	_kind string `boom:"kind,User"`
	ID    string `boom:"id"`
	Role  int32
}

func onlyID(id domain.UserID) *entity {
	return &entity{ID: id.String()}
}

func (e *entity) toDomain() *domain.User {
	return &domain.User{
		ID:   domain.UserID(e.ID),
		Role: pb.User_Role(e.Role),
	}
}

func toEntity(from *domain.User) *entity {
	return &entity{
		ID:   from.ID.String(),
		Role: int32(from.Role),
	}
}

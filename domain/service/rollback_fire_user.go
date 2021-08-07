package service

import (
	"context"

	"gae-go-recruiting-server/adapter"
	"gae-go-recruiting-server/domain"
)

func NewRollbackFireUserService(
	userRepository adapter.UserRepository,
	fireUserRepository adapter.FireUserRepository,
	logger adapter.CompositeLogger) adapter.RollbackFireUserService {
	return func(ctx context.Context, userID domain.UserID) {
		if isExists, _ := userRepository.Exists(ctx, userID); isExists {
			return
		}

		if err := fireUserRepository.Delete(ctx, userID); err != nil {
			logger.Error().With(ctx).Printf("cannot delete fire user of %s when failed customer create", userID)
		}
	}
}

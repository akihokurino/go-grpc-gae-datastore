package prefecture_app

import (
	"context"

	"gae-go-recruiting-server/adapter"
	"gae-go-recruiting-server/domain"

	"github.com/pkg/errors"
)

type application struct {
	logger               adapter.CompositeLogger
	prefectureRepository adapter.PrefectureRepository
}

func NewApplication(
	logger adapter.CompositeLogger,
	prefectureRepository adapter.PrefectureRepository) adapter.PrefectureApplication {
	return &application{
		logger:               logger,
		prefectureRepository: prefectureRepository,
	}
}

func (a *application) GetAll(ctx context.Context) ([]*domain.Prefecture, error) {
	prefectures, err := a.prefectureRepository.GetAll(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return prefectures, nil
}

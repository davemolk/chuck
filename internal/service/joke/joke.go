package joke

import (
	"context"

	"github.com/davemolk/chuck/internal/domain"
	"github.com/davemolk/chuck/internal/service"
	"go.uber.org/zap"
)

var _ service.JokeService = (*Service)(nil)

type chuckGetter interface {
	GetRandomJoke(ctx context.Context) (*domain.Joke, error)
}

type Service struct {
	logger *zap.Logger
	client chuckGetter
}

func NewService(logger *zap.Logger, client chuckGetter) *Service {
	return &Service{
		logger: logger,
		client: client,
	}
}

func (s *Service) GetJoke(ctx context.Context) (*domain.Joke, error) {
	return s.client.GetRandomJoke(ctx)
}

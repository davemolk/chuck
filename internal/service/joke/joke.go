package joke

import (
	"context"
	"fmt"

	"github.com/davemolk/chuck/internal/domain"
	"github.com/davemolk/chuck/internal/service"
	"go.uber.org/zap"
)

const maxJokesFromAPI = 30

var _ service.JokeService = (*Service)(nil)

type chuckGetter interface {
	Search(ctx context.Context, query string, limit int) ([]*domain.Joke, error)
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
	// return s.client.GetRandomJoke(ctx)
	return nil, nil
}

func (s *Service) GetRandomJoke(ctx context.Context, query string) (*domain.Joke, error) {
	// check database

	// call api
	logger := s.logger.With(zap.String("query", query))
	logger.Info("no cached matches, calling api...")

	jokes, err := s.client.Search(ctx, query, maxJokesFromAPI)
	if err != nil {
		return nil, fmt.Errorf("failed to search api: %w", err)
	}

	if len(jokes) == 0 {
		logger.Info("no jokes found from api")
		return nil, nil
	}

	logger.Info("found", zap.Int("count", len(jokes)))

	// populate db

	// return
	return jokes[0], nil
}

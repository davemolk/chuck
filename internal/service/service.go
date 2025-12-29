package service

import (
	"context"

	"github.com/davemolk/chuck/internal/domain"
)

type JokeService interface {
	GetPersonalizedJoke(ctx context.Context, name string) (*domain.Joke, error)
	GetRandomJokeByQuery(ctx context.Context, query string) (*domain.Joke, error)
}

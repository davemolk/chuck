package service

import (
	"context"

	"github.com/davemolk/chuck/internal/domain"
)

type JokeService interface {
	GetJoke(ctx context.Context) (*domain.Joke, error)
	GetRandomJoke(ctx context.Context, query string) (*domain.Joke, error)
}

package service

import (
	"context"
	"time"

	"github.com/davemolk/chuck/internal/domain"
)

type JokeService interface {
	GetPersonalizedJoke(ctx context.Context, name string) (*domain.Joke, error)
	GetRandomJokeByQuery(ctx context.Context, query string) (*domain.Joke, error)
}

type TokenService interface {
	CreateToken(ctx context.Context, userID int64, ttl time.Duration) (*domain.Token, error)
}

type UserService interface {
	CreateUser(ctx context.Context, email, password string) (int64, error)
}

package service

import (
	"context"
	"time"

	"github.com/davemolk/chuck/internal/domain"
)

type JokeService interface {
	GetPersonalizedJoke(ctx context.Context, name string) (*domain.Joke, error)
	GetRandomJoke(ctx context.Context) (*domain.Joke, error)
	GetRandomJokeByQuery(ctx context.Context, query string) (*domain.Joke, error)
}

type TokenService interface {
	CreateToken(ctx context.Context, userID int64, ttl time.Duration) (*domain.Token, error)
	ValidateToken(ctx context.Context, token string) (int64, error)
}

type UserService interface {
	CreateUser(ctx context.Context, email, password string) (int64, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	GetUserByID(ctx context.Context, id int64) (*domain.User, error)
}

type AuthService interface {
	Login(ctx context.Context, email, password string) (*domain.Token, error)
	GetUserIDForToken(ctx context.Context, token string) (*domain.User, error)
}

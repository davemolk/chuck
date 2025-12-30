package mock

import (
	"context"
	"time"

	"github.com/davemolk/chuck/internal/domain"
)

type ChuckClient struct {
	SearchFn     func(ctx context.Context, query string, limit int) ([]*domain.Joke, error)
	SearchCalled bool
}

func (c *ChuckClient) Search(ctx context.Context, query string, limit int) ([]*domain.Joke, error) {
	c.SearchCalled = true
	return c.SearchFn(ctx, query, limit)
}

type UserService struct {
	CreateUserFn         func(ctx context.Context, email, password string) (int64, error)
	CreateUserCalled     bool
	GetUserByEmailFn     func(ctx context.Context, email string) (*domain.User, error)
	GetUserByEmailCalled bool
	GetUserByIDFn        func(ctx context.Context, id int64) (*domain.User, error)
	GetUserByIDCalled    bool
}

func (s *UserService) CreateUser(ctx context.Context, email, password string) (int64, error) {
	s.CreateUserCalled = true
	return s.CreateUserFn(ctx, email, password)
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	s.GetUserByEmailCalled = true
	return s.GetUserByEmailFn(ctx, email)
}

func (s *UserService) GetUserByID(ctx context.Context, id int64) (*domain.User, error) {
	s.GetUserByIDCalled = true
	return s.GetUserByIDFn(ctx, id)
}

type TokenService struct {
	CreateTokenFn       func(ctx context.Context, userID int64, ttl time.Duration) (*domain.Token, error)
	CreateTokenFnCalled bool
	ValidateTokenFn     func(ctx context.Context, token string) (int64, error)
	ValidateTokenCalled bool
}

func (s *TokenService) CreateToken(ctx context.Context, userID int64, ttl time.Duration) (*domain.Token, error) {
	s.CreateTokenFnCalled = true
	return s.CreateTokenFn(ctx, userID, ttl)
}

func (s *TokenService) ValidateToken(ctx context.Context, token string) (int64, error) {
	s.ValidateTokenCalled = true
	return s.ValidateTokenFn(ctx, token)
}

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

func (s *UserService) ResetCalls() {
	s.CreateUserCalled = false
	s.GetUserByIDCalled = false
	s.GetUserByEmailCalled = false
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

type JokeService struct {
	GetPersonalizedJokeFn      func(ctx context.Context, name string) (*domain.Joke, error)
	GetPersonalizedJokeCalled  bool
	GetRandomJokeFn            func(ctx context.Context) (*domain.Joke, error)
	GetRandomJokeCalled        bool
	GetRandomJokeByQueryFn     func(ctx context.Context, query string) (*domain.Joke, error)
	GetRandomJokeByQueryCalled bool
}

func (s *JokeService) GetPersonalizedJoke(ctx context.Context, name string) (*domain.Joke, error) {
	s.GetPersonalizedJokeCalled = true
	return s.GetPersonalizedJokeFn(ctx, name)
}

func (s *JokeService) GetRandomJoke(ctx context.Context) (*domain.Joke, error) {
	s.GetRandomJokeCalled = true
	return s.GetRandomJokeFn(ctx)
}

func (s *JokeService) GetRandomJokeByQuery(ctx context.Context, query string) (*domain.Joke, error) {
	s.GetRandomJokeByQueryCalled = true
	return s.GetRandomJokeByQueryFn(ctx, query)
}

func (s *JokeService) ResetCalls() {
	s.GetPersonalizedJokeCalled = false
	s.GetRandomJokeByQueryCalled = false
	s.GetRandomJokeCalled = false
}

type AuthService struct {
	LoginFn                 func(ctx context.Context, email, password string) (*domain.Token, error)
	LoginFnCalled           bool
	GetUserIDForTokenFn     func(ctx context.Context, token string) (*domain.User, error)
	GetUserIDForTokenCalled bool
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*domain.Token, error) {
	s.LoginFnCalled = true
	return s.LoginFn(ctx, email, password)
}

func (s *AuthService) GetUserIDForToken(ctx context.Context, token string) (*domain.User, error) {
	s.GetUserIDForTokenCalled = true
	return s.GetUserIDForTokenFn(ctx, token)
}

func (s *AuthService) ResetCalls() {
	s.LoginFnCalled = false
	s.GetUserIDForTokenCalled = false
}

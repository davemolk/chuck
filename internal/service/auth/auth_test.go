package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/davemolk/chuck/internal/domain"
	"github.com/davemolk/chuck/internal/service/token"
	"github.com/davemolk/chuck/internal/service/user"
	"github.com/davemolk/chuck/internal/sql/dbtest"
	"github.com/davemolk/chuck/internal/tests/fixture"
	"github.com/davemolk/chuck/internal/tests/mock"
	"github.com/stretchr/testify/require"
)

func TestGetUserForToken(t *testing.T) {
	db := dbtest.SetupTestDB(t)
	email := "roundhouse@kick.com"
	userService := &mock.UserService{
		GetUserByIDFn: func(ctx context.Context, id int64) (*domain.User, error) {
			return &domain.User{
				Email: email,
			}, nil
		},
	}
	tokenService := &mock.TokenService{
		ValidateTokenFn: func(ctx context.Context, token string) (int64, error) {
			return 1, nil
		},
	}

	s := NewService(fixture.TestLogger(t), db, userService, tokenService)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		user, err := s.GetUserIDForToken(ctx, "my token")
		require.NoError(t, err)
		require.Equal(t, email, user.Email)
	})

	t.Run("error: invalid token", func(t *testing.T) {
		invalidTokenService := &mock.TokenService{
			ValidateTokenFn: func(ctx context.Context, tokenStr string) (int64, error) {
				return 0, token.ErrInvalidToken
			},
		}
		s.tokenService = invalidTokenService
		_, err := s.GetUserIDForToken(ctx, "my token")
		require.Error(t, err)
		require.True(t, errors.Is(err, token.ErrInvalidToken))
		s.tokenService = tokenService
	})

	t.Run("error: invalid user", func(t *testing.T) {
		invalidUserService := &mock.UserService{
			GetUserByIDFn: func(ctx context.Context, id int64) (*domain.User, error) {
				return nil, errors.New("nope")
			},
		}
		s.userService = invalidUserService
		_, err := s.GetUserIDForToken(ctx, "my token")
		require.Error(t, err)
	})
}

func TestLogin(t *testing.T) {
	db := dbtest.SetupTestDB(t)
	email := "roundhouse@kick.com"
	pw := "pw"
	ctx := context.Background()

	// use real service so we get proper hashed password
	userService := user.NewService(fixture.TestLogger(t), db)
	userID, err := userService.CreateUser(ctx, email, pw)
	require.NoError(t, err)

	tokenService := &mock.TokenService{
		CreateTokenFn: func(ctx context.Context, userID int64, ttl time.Duration) (*domain.Token, error) {
			return &domain.Token{
				UserID:    userID,
				Plaintext: "roundhouse",
			}, nil
		},
	}

	s := NewService(fixture.TestLogger(t), db, userService, tokenService)

	t.Run("success", func(t *testing.T) {
		token, err := s.Login(ctx, email, pw)
		require.NoError(t, err)
		require.Equal(t, "roundhouse", token.Plaintext)
		require.Equal(t, userID, token.UserID)
	})

	t.Run("error: no user", func(t *testing.T) {
		_, err := s.Login(ctx, "no@email.com", pw)
		require.Error(t, err)
		require.True(t, errors.Is(err, ErrInvalidCredentials))
	})
}

package user

import (
	"context"
	"errors"
	"testing"

	"github.com/davemolk/chuck/internal/domain"
	"github.com/davemolk/chuck/internal/sql/dbtest"
	"github.com/davemolk/chuck/internal/tests/fixture"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	db := dbtest.SetupTestDB(t)
	s := NewService(fixture.TestLogger(t), db)
	ctx := context.Background()

	email := "chuck@norris.com"
	pw := "r0undhou5e"
	t.Run("success", func(t *testing.T) {
		id, err := s.CreateUser(ctx, email, pw)
		require.NoError(t, err)
		require.NotEqual(t, 0, id)
	})

	t.Run("can't use same email", func(t *testing.T) {
		newPW := "foo"
		_, err := s.CreateUser(ctx, email, newPW)
		require.Error(t, err)
		require.True(t, errors.Is(err, ErrDuplicateEmail))
	})
}

func TestGetUserByEmail(t *testing.T) {
	db := dbtest.SetupTestDB(t)
	s := NewService(fixture.TestLogger(t), db)
	ctx := context.Background()

	email := "chuck@norris.com"
	t.Run("error: user not exist", func(t *testing.T) {
		_, err := s.GetUserByEmail(ctx, email)
		require.Error(t, err)
		require.True(t, errors.Is(err, domain.ErrNotFound))
	})

	pw := "r0undhou5e"
	id, err := s.CreateUser(ctx, email, pw)
	require.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		user, err := s.GetUserByEmail(ctx, email)
		require.NoError(t, err)
		require.Equal(t, id, user.ID)
	})
}

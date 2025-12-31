package token

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/davemolk/chuck/internal/sql/dbtest"
	"github.com/davemolk/chuck/internal/tests/fixture"
	"github.com/stretchr/testify/require"
)

func TestCreateToken(t *testing.T) {
	db := dbtest.SetupTestDB(t)
	s := NewService(fixture.TestLogger(t), db)
	ctx := context.Background()
	ttl := 5 * time.Minute
	t.Run("error: no user", func(t *testing.T) {
		_, err := s.CreateToken(ctx, 20, ttl)
		require.Error(t, err)
	})

	user1ID := fixture.AddUser(t, db, "email1")
	user2ID := fixture.AddUser(t, db, "email2")
	t.Run("create two tokens", func(t *testing.T) {
		now := time.Now()
		token1, err := s.CreateToken(ctx, user1ID, ttl)
		require.NoError(t, err)

		token2, err := s.CreateToken(ctx, user2ID, ttl)
		require.NoError(t, err)

		require.NotEqual(t, token1.Plaintext, token2.Plaintext)
		require.NotEqual(t, token1.UserID, token2.UserID)
		require.NotEqual(t, token1.Hash, token2.Hash)

		require.Equal(t, user1ID, token1.UserID)
		require.Equal(t, user2ID, token2.UserID)

		require.Greater(t, now.Add(10*time.Minute), token1.ExpiresAt)
		require.Less(t, now.Add(1*time.Minute), token1.ExpiresAt)
	})
}

func TestValidateToken(t *testing.T) {
	db := dbtest.SetupTestDB(t)
	s := NewService(fixture.TestLogger(t), db)
	ctx := context.Background()
	ttl := 5 * time.Minute

	t.Run("error: no user", func(t *testing.T) {
		_, err := s.ValidateToken(ctx, "foobar")
		require.Error(t, err)
		require.True(t, errors.Is(err, ErrInvalidToken))
	})

	user1ID := fixture.AddUser(t, db, "email1")
	user2ID := fixture.AddUser(t, db, "email2")
	token1, err := s.CreateToken(ctx, user1ID, ttl)
	require.NoError(t, err)

	token2, err := s.CreateToken(ctx, user2ID, ttl)
	require.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		id, err := s.ValidateToken(ctx, token1.Plaintext)
		require.NoError(t, err)
		require.Equal(t, user1ID, id)

		// confirm other
		id2, err := s.ValidateToken(ctx, token2.Plaintext)
		require.NoError(t, err)
		require.Equal(t, user2ID, id2)
	})
}

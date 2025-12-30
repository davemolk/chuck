package fixture

import (
	"context"
	"testing"

	"github.com/davemolk/chuck/internal/sql"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func TestLogger(t *testing.T) *zap.Logger {
	logger, err := zap.NewDevelopment()
	require.NoError(t, err)
	return logger
}

func AddUser(t *testing.T, db *sql.DB, email string) int64 {
	password := "password"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	require.NoError(t, err)

	args := []any{email, hash}

	query := `
		insert into users (email, hashed_pw)
		values ($1, $2)
		on conflict (email) do nothing
		returning id`

	var id int64
	err = db.QueryRowContext(context.Background(), query, args...).Scan(&id)
	require.NoError(t, err)

	return id
}

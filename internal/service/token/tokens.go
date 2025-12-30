package token

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base32"
	"errors"
	"fmt"
	"time"

	"github.com/davemolk/chuck/internal/domain"
	"github.com/davemolk/chuck/internal/service"
	sqldb "github.com/davemolk/chuck/internal/sql"
	"go.uber.org/zap"
)

var ErrInvalidToken = errors.New("invalid token")

var _ service.TokenService = (*Service)(nil)

type Service struct {
	logger *zap.Logger
	db     *sqldb.DB
}

func NewService(logger *zap.Logger, db *sqldb.DB) *Service {
	return &Service{
		logger: logger,
		db:     db,
	}
}

func (s *Service) generateToken(userID int64, ttl time.Duration) (*domain.Token, error) {
	token := &domain.Token{
		UserID:    userID,
		ExpiresAt: time.Now().Add(ttl),
	}

	randomBytes := make([]byte, 16)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]

	return token, nil
}

func (s *Service) CreateToken(ctx context.Context, userID int64, ttl time.Duration) (*domain.Token, error) {
	token, err := s.generateToken(userID, ttl)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	query := `
	insert into tokens (hash, user_id, expires_at)
	values ($1, $2, $3)
	`

	args := []any{token.Hash, token.UserID, token.ExpiresAt}

	if _, err = s.db.ExecContext(ctx, query, args...); err != nil {
		return nil, fmt.Errorf("failed to insert token: %w", err)
	}

	return token, nil
}

func (s *Service) ValidateToken(ctx context.Context, token string) (int64, error) {
	hash := sha256.Sum256([]byte(token))

	query := `select user_id from tokens where hash = $1 and expires_at > $2`

	args := []any{hash[:], time.Now()}

	var id int64
	if err := s.db.QueryRowContext(ctx, query, args...).Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidToken
		}
		return 0, fmt.Errorf("failed to validate token: %w", err)
	}

	return id, nil
}

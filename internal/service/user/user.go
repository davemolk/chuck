package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/davemolk/chuck/internal/domain"
	sqldb "github.com/davemolk/chuck/internal/sql"
	"golang.org/x/crypto/bcrypt"

	"github.com/davemolk/chuck/internal/service"
	"go.uber.org/zap"
)

var ErrDuplicateEmail = errors.New("duplicate email")

type Service struct {
	logger       *zap.Logger
	db           *sqldb.DB
	tokenService service.TokenService
}

var _ service.UserService = (*Service)(nil)

func NewService(logger *zap.Logger, db *sqldb.DB, tokenService service.TokenService) *Service {
	return &Service{
		logger:       logger,
		db:           db,
		tokenService: tokenService,
	}
}

func (s *Service) CreateUser(ctx context.Context, email, password string) (int64, error) {
	logger := s.logger.With(zap.String("email", email))
	logger.Info("creating user")

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return 0, fmt.Errorf("failed to hash: %w", err)
	}

	args := []any{email, hash}

	query := `
		insert into users (email, hashed_pw)
		values ($1, $2)
		on conflict (email) do nothing
		returning id`

	var id int64
	err = s.db.QueryRowContext(ctx, query, args...).Scan(&id)
	if err != nil {
		// this would happen if the email is already stored, so
		// we do nothing and consequently can't scan the id
		if errors.Is(err, sql.ErrNoRows) {
			logger.Debug("email already exists")
			return 0, ErrDuplicateEmail
		}
		return 0, fmt.Errorf("failed to insert user: %w", err)
	}

	logger.Info("user created")

	return id, nil
}

func (s *Service) getByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `select id, email, hashed_pw, created_at from users where email = $1`

	var u domain.User
	err := s.db.QueryRowContext(ctx, query, email).Scan(
		&u.ID,
		&u.Email,
		&u.HashedPW,
		&u.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &u, nil
}

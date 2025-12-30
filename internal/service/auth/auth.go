package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/davemolk/chuck/internal/domain"
	"github.com/davemolk/chuck/internal/service"
	sqldb "github.com/davemolk/chuck/internal/sql"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var _ service.AuthService = (*Service)(nil)

var ErrInvalidCredentials = errors.New("invalid credentials")

type Service struct {
	logger       *zap.Logger
	db           *sqldb.DB
	userService  service.UserService
	tokenService service.TokenService
}

func NewService(logger *zap.Logger, db *sqldb.DB, userService service.UserService, tokenService service.TokenService) *Service {
	return &Service{
		logger:       logger,
		db:           db,
		userService:  userService,
		tokenService: tokenService,
	}
}

func (s *Service) GetUserForToken(ctx context.Context, token string) (*domain.User, error) {
	s.logger.Info("validating token", zap.String("t", token))
	userID, err := s.tokenService.ValidateToken(ctx, token)
	if err != nil {
		return nil, err
	}

	s.logger.Info("getting user")
	user, err := s.userService.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			// valid token but no user, should invalidate token
			return nil, errors.New("no user found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (s *Service) Login(ctx context.Context, email, password string) (*domain.Token, error) {
	user, err := s.userService.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	valid, err := s.validatePasswordHash(user.HashedPW, password)
	if err != nil {
		// log these attempts?
		return nil, ErrInvalidCredentials
	}

	if !valid {
		return nil, ErrInvalidCredentials
	}

	token, err := s.tokenService.CreateToken(ctx, user.ID, 24*time.Hour)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (s *Service) validatePasswordHash(hash []byte, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/davemolk/chuck/internal/domain"
)

type contextKey string

const requestIDKey contextKey = "request_id"
const requestIDHeader = "X-Request-ID"
const userKey contextKey = "user"

func generateRequestID() string {
	b := make([]byte, 6)

	// docstring says it never returns an error, safe to ignore
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func RequestIDFromCtx(ctx context.Context) string {
	if id, ok := ctx.Value(requestIDKey).(string); ok {
		return id
	}

	return ""
}

func UserFromCtx(ctx context.Context) (*domain.User, error) {
	user, ok := ctx.Value(userKey).(*domain.User)
	if !ok {
		return nil, fmt.Errorf("missing user")
	}

	return user, nil
}

func UserToCtx(ctx context.Context, user *domain.User) context.Context {
	ctx = context.WithValue(ctx, userKey, user)
	return ctx
}

package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
)

type contextKey string

const requestIDKey contextKey = "request_id"
const requestIDHeader = "X-Request-ID"

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

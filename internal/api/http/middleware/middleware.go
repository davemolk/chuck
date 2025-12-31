package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/davemolk/chuck/internal/domain"
	"go.uber.org/zap"
)

type respWriter struct {
	wrapped       http.ResponseWriter
	statusCode    int
	headerWritten bool
}

func newRespWriter(w http.ResponseWriter) *respWriter {
	return &respWriter{
		wrapped:    w,
		statusCode: http.StatusOK,
	}
}

func (rw *respWriter) Header() http.Header {
	return rw.wrapped.Header()
}

func (rw *respWriter) WriteHeader(statusCode int) {
	rw.wrapped.WriteHeader(statusCode)

	if !rw.headerWritten {
		rw.statusCode = statusCode
		rw.headerWritten = true
	}
}

func (rw *respWriter) Write(b []byte) (int, error) {
	rw.headerWritten = true
	return rw.wrapped.Write(b)
}

func Logger(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			wrapped := newRespWriter(w)
			reqID := RequestIDFromCtx(r.Context())

			var email string
			user, err := UserFromCtx(r.Context())
			// todo: come back to this
			if err == nil && user != nil {
				email = user.Email
			}

			logger.Info("request started", zap.String("request_id", reqID), zap.String("email", email), zap.String("method", r.Method), zap.String("path", r.URL.Path), zap.String("ua", r.UserAgent()))

			next.ServeHTTP(wrapped, r)

			duration := time.Since(start)
			logger.Info("request finished", zap.String("request_id", reqID), zap.String("email", email), zap.String("method", r.Method), zap.String("path", r.URL.Path), zap.String("ua", r.UserAgent()), zap.Int("status", wrapped.statusCode), zap.Duration("duration", duration))
		})
	}
}

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get(requestIDHeader)
		if requestID == "" {
			requestID = generateRequestID()
		}

		ctx := context.WithValue(r.Context(), requestIDKey, requestID)
		w.Header().Set(requestIDHeader, requestID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RecoverPanic(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					requestID := RequestIDFromCtx(r.Context())
					logger.Error("panic recovered", zap.String("request_id", requestID), zap.Any("error", err), zap.String("path", r.URL.Path), zap.String("stack", string(debug.Stack())))

					w.Header().Set("Connection", "close")
					http.Error(w, fmt.Sprintf("internal server error: requestID: %s", requestID), http.StatusInternalServerError)
					return
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

type userAuthenticator interface {
	GetUserIDForToken(ctx context.Context, token string) (*domain.User, error)
}

func Auth(userAuthenticator userAuthenticator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				next.ServeHTTP(w, r)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "invalid or missing authentication token", http.StatusUnauthorized)
				return
			}

			token := parts[1]

			user, err := userAuthenticator.GetUserIDForToken(r.Context(), token)
			if err != nil {
				if errors.Is(err, domain.ErrNotFound) {
					http.Error(w, "invalid token", http.StatusUnauthorized)
				} else {
					http.Error(w, "server is unable to process request", http.StatusInternalServerError)
				}
				return
			}

			userCtx := UserToCtx(r.Context(), user)
			r = r.WithContext(userCtx)
			next.ServeHTTP(w, r)
		})
	}
}

func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, err := UserFromCtx(r.Context()); err != nil {
			http.Error(w, "authentication required", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}

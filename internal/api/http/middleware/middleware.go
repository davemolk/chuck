package middleware

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

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

			logger.Info("request started", zap.String("request_id", reqID), zap.String("method", r.Method), zap.String("path", r.URL.Path), zap.String("ua", r.UserAgent()))

			next.ServeHTTP(wrapped, r)
			duration := time.Since(start)

			logger.Info("request finished", zap.String("request_id", reqID), zap.String("method", r.Method), zap.String("path", r.URL.Path), zap.String("ua", r.UserAgent()), zap.Int("status", wrapped.statusCode), zap.Duration("duration", duration))
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

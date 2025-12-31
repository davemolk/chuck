package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/davemolk/chuck/internal/domain"
	"github.com/davemolk/chuck/internal/service/auth"
	"github.com/davemolk/chuck/internal/service/joke"
	"github.com/davemolk/chuck/internal/service/token"
	"github.com/davemolk/chuck/internal/service/user"

	"github.com/davemolk/chuck/internal/api/http/middleware"
	"go.uber.org/zap"
)

func respondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			return
		}
	}
}

func readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return json.NewDecoder(r.Body).Decode(&data)
}

type errResponse struct {
	Error      string
	RequestID  string
	StatusCode int
}

func respondError(w http.ResponseWriter, r *http.Request, logger *zap.Logger, status int, err error) {
	requestID := middleware.RequestIDFromCtx(r.Context())
	logger.Error("request error", zap.Int("status", status), zap.String("request_id", requestID), zap.Error(err))
	respondJSON(w, status, errResponse{
		Error:      err.Error(),
		RequestID:  requestID,
		StatusCode: status,
	})
}

func errToStatusCode(err error) int {
	switch err {
	case domain.ErrNotFound:
		return http.StatusNotFound
	case joke.ErrNoJokes:
		return http.StatusNotFound
	case auth.ErrInvalidCredentials:
		return http.StatusUnauthorized
	case token.ErrInvalidToken:
		return http.StatusUnauthorized
	case user.ErrDuplicateEmail:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

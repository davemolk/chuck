package handlers

import (
	"errors"
	"net/http"

	"github.com/davemolk/chuck/internal/service"
	"go.uber.org/zap"
)

type AuthHandlers struct {
	logger      *zap.Logger
	authService service.AuthService
}

func NewAuthHandlers(logger *zap.Logger, authService service.AuthService) *AuthHandlers {
	return &AuthHandlers{
		logger:      logger,
		authService: authService,
	}
}

func (h *AuthHandlers) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := readJSON(w, r, &req); err != nil {
		respondError(w, r, h.logger, http.StatusBadRequest, err)
		return
	}

	if req.Email == "" {
		respondError(w, r, h.logger, http.StatusBadRequest, errors.New("email required"))
		return
	}

	if req.Password == "" {
		respondError(w, r, h.logger, http.StatusBadRequest, errors.New("password required"))
		return
	}

	token, err := h.authService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		respondError(w, r, h.logger, errToStatusCode(err), err)
		return
	}

	respondJSON(w, http.StatusOK, token)
}

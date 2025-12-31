package handlers

import (
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

	if err := validateEmail(req.Email); err != nil {
		respondError(w, r, h.logger, http.StatusBadRequest, err)
		return
	}

	if err := validatePassword(req.Password); err != nil {
		respondError(w, r, h.logger, http.StatusBadRequest, err)
		return
	}

	token, err := h.authService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		respondError(w, r, h.logger, errToStatusCode(err), err)
		return
	}

	respondJSON(w, http.StatusOK, token)
}

package handlers

import (
	"net/http"

	"github.com/davemolk/chuck/internal/service"
	"go.uber.org/zap"
)

type UserHandlers struct {
	logger      *zap.Logger
	userService service.UserService
}

func NewUserHandlers(logger *zap.Logger, userService service.UserService) *UserHandlers {
	return &UserHandlers{
		logger:      logger,
		userService: userService,
	}
}

func (h *UserHandlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := readJSON(w, r, &req); err != nil {
		respondError(w, r, h.logger, errToStatusCode(err), err)
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

	user, err := h.userService.CreateUser(r.Context(), req.Email, req.Password)
	if err != nil {
		respondError(w, r, h.logger, errToStatusCode(err), err)
		return
	}

	respondJSON(w, http.StatusCreated, user)
}

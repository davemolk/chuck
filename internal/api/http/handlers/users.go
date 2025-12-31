package handlers

import (
	"errors"
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

	if req.Email == "" {
		respondError(w, r, h.logger, http.StatusBadRequest, errors.New("email required"))
		return
	}

	if req.Password == "" {
		respondError(w, r, h.logger, http.StatusBadRequest, errors.New("password required"))
		return
	}

	user, err := h.userService.CreateUser(r.Context(), req.Email, req.Password)
	if err != nil {
		respondError(w, r, h.logger, errToStatusCode(err), err)
		return
	}

	respondJSON(w, http.StatusCreated, user)
}

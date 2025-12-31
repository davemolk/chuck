package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/davemolk/chuck/internal/service"
	"go.uber.org/zap"
)

// minQueryLength is a limit set by the Chuck Norris API.
const minQueryLength = 3

type JokeHandlers struct {
	logger      *zap.Logger
	jokeService service.JokeService
}

func NewJokeHandlers(logger *zap.Logger, jokeService service.JokeService) *JokeHandlers {
	return &JokeHandlers{
		logger:      logger,
		jokeService: jokeService,
	}
}

func (h *JokeHandlers) GetPersonalizedJoke(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		respondError(w, r, h.logger, http.StatusBadRequest, errors.New("name is required"))
		return
	}

	joke, err := h.jokeService.GetPersonalizedJoke(r.Context(), name)
	if err != nil {
		respondError(w, r, h.logger, errToStatusCode(err), err)
		return
	}

	respondJSON(w, http.StatusOK, joke)
}

func (h *JokeHandlers) GetRandomJoke(w http.ResponseWriter, r *http.Request) {
	joke, err := h.jokeService.GetRandomJoke(r.Context())
	if err != nil {
		respondError(w, r, h.logger, errToStatusCode(err), err)
		return
	}

	respondJSON(w, http.StatusOK, joke)
}

func (h *JokeHandlers) GetRandomJokeByQuery(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	if len(query) < minQueryLength {
		respondError(w, r, h.logger, http.StatusBadRequest, fmt.Errorf("query of minimum length %d is required", minQueryLength))
		return
	}

	joke, err := h.jokeService.GetRandomJokeByQuery(r.Context(), query)
	if err != nil {
		respondError(w, r, h.logger, errToStatusCode(err), err)
		return
	}

	respondJSON(w, http.StatusOK, joke)
}

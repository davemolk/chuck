package handlers

import (
	"net/http"

	"github.com/davemolk/chuck/internal/service"
	"go.uber.org/zap"
)

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

	if err := validateName(name); err != nil {
		respondError(w, r, h.logger, http.StatusBadRequest, err)
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

	if err := validateQuery(query); err != nil {
		respondError(w, r, h.logger, http.StatusBadRequest, err)
		return
	}

	joke, err := h.jokeService.GetRandomJokeByQuery(r.Context(), query)
	if err != nil {
		respondError(w, r, h.logger, errToStatusCode(err), err)
		return
	}

	respondJSON(w, http.StatusOK, joke)
}

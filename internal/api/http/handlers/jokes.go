package handlers

import (
	"fmt"
	"net/http"

	"github.com/davemolk/chuck/internal/service"
	"go.uber.org/zap"
)

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

func (h *JokeHandlers) GetPersonalized(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if len(name) == 0 {
		// todo: handle this
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	joke, err := h.jokeService.GetPersonalizedJoke(r.Context(), name)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, joke)
}

func (h *JokeHandlers) GetRandom(w http.ResponseWriter, r *http.Request) {
	joke, err := h.jokeService.GetRandomJoke(r.Context())
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, joke)
}

func (h *JokeHandlers) GetRandomByQuery(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	if len(query) < minQueryLength {
		http.Error(w, fmt.Sprintf("query of minimum length %d is required", minQueryLength), http.StatusBadRequest)
		return
	}

	joke, err := h.jokeService.GetRandomJokeByQuery(r.Context(), query)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, joke)
}

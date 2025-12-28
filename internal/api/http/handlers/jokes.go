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

func (h *JokeHandlers) GetRandom(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	if len(query) < minQueryLength {
		http.Error(w, fmt.Sprintf("query of minimum length %d is required", minQueryLength), http.StatusBadRequest)
		return
	}

	joke, err := h.jokeService.GetRandomJoke(r.Context(), query)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, joke)
}

func (h *JokeHandlers) GetPersonalized(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
	}

	if err := readJSON(w, r, &req); err != nil {
		// todo
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusBadRequest)
		return
	}

	if len(req.Name) == 0 {
		// todo: handle this
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	joke, err := h.jokeService.GetJoke(r.Context())
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	// todo: ignore for now, this will be swapped for a reusable writer helper
	_, _ = w.Write([]byte(joke.Content))
}

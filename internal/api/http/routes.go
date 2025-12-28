package http

import (
	"net/http"

	"github.com/davemolk/chuck/internal/api/http/handlers"
	"github.com/davemolk/chuck/internal/service"
	"go.uber.org/zap"
)

type Services struct {
	JokeService service.JokeService
}

func NewRoutes(logger *zap.Logger, services *Services) http.Handler {
	mux := http.NewServeMux()
	jokes := handlers.NewJokeHandlers(logger, services.JokeService)

	mux.HandleFunc("POST /api/v1/jokes/personalized", jokes.GetPersonalized)
	mux.HandleFunc("GET /api/v1/jokes/random", jokes.GetRandom)

	return mux
}

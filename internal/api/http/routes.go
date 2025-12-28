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

	mux.HandleFunc("GET /v1/jokes", jokes.GetRandom)

	return mux
}

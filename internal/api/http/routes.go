package http

import (
	"net/http"

	"github.com/davemolk/chuck/internal/api/http/handlers"
	"github.com/davemolk/chuck/internal/api/http/middleware"
	"github.com/davemolk/chuck/internal/service"
	"go.uber.org/zap"
)

type Services struct {
	JokeService service.JokeService
	UserService service.UserService
}

func NewRoutes(logger *zap.Logger, services *Services) http.Handler {
	mux := http.NewServeMux()

	health := handlers.NewHealthHandlers()
	jokes := handlers.NewJokeHandlers(logger, services.JokeService)
	users := handlers.NewUserHandlers(logger, services.UserService)

	mux.HandleFunc("GET /health", health.HealthCheck)
	mux.HandleFunc("POST /api/v1/jokes/personalized", jokes.GetPersonalized)
	mux.HandleFunc("GET /api/v1/jokes/random", jokes.GetRandom)
	mux.HandleFunc("POST /api/v1/users", users.CreateUser)

	var handler http.Handler = mux
	handler = middleware.Logger(logger)(handler)
	handler = middleware.RequestID(handler)
	handler = middleware.RecoverPanic(logger)(handler)

	return handler
}

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
	AuthService service.AuthService
}

func NewRoutes(logger *zap.Logger, services *Services) http.Handler {
	mux := http.NewServeMux()

	health := handlers.NewHealthHandlers()
	jokes := handlers.NewJokeHandlers(logger, services.JokeService)
	users := handlers.NewUserHandlers(logger, services.UserService)
	auth := handlers.NewAuthHandlers(logger, services.AuthService)

	mux.HandleFunc("GET /health", health.HealthCheck)

	mux.HandleFunc("GET /api/v1/jokes/random", jokes.GetRandom)
	mux.HandleFunc("GET /api/v1/jokes/search", middleware.RequireAuth(jokes.GetRandomByQuery))
	mux.HandleFunc("GET /api/v1/jokes/personalized", middleware.RequireAuth(jokes.GetPersonalized))

	mux.HandleFunc("POST /api/v1/users", users.CreateUser)
	mux.HandleFunc("POST /api/v1/auth/login", auth.Login)

	var handler http.Handler = mux
	handler = middleware.Logger(logger)(handler)
	handler = middleware.Auth(services.AuthService)(handler)
	handler = middleware.RequestID(handler)
	handler = middleware.RecoverPanic(logger)(handler)

	return handler
}

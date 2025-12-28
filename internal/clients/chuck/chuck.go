package chuck

import (
	"context"
	"net/http"

	"github.com/davemolk/chuck/internal/domain"
	"go.uber.org/zap"
)

const randomJokeURL = "https://api.chucknorris.io/jokes/random"

type APIClient struct {
	Logger  *zap.Logger
	Client  *http.Client
	baseURL string
}

func NewClient(logger *zap.Logger) *APIClient {
	return &APIClient{
		Logger:  logger,
		Client:  http.DefaultClient,
		baseURL: randomJokeURL,
	}
}

func (c *APIClient) GetRandomJoke(ctx context.Context) (*domain.Joke, error) {
	return &domain.Joke{
		Content: "round house!",
	}, nil
}

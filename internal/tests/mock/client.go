package mock

import (
	"context"

	"github.com/davemolk/chuck/internal/domain"
)

type ChuckClient struct {
	SearchFn     func(ctx context.Context, query string, limit int) ([]*domain.Joke, error)
	SearchCalled bool
}

func (c *ChuckClient) Search(ctx context.Context, query string, limit int) ([]*domain.Joke, error) {
	c.SearchCalled = true
	return c.SearchFn(ctx, query, limit)
}

package chuck

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/davemolk/chuck/internal/domain"
	"go.uber.org/zap"
)

const (
	baseURL            = "https://api.chucknorris.io"
	responseTimeFormat = "2006-01-02 15:04:05"
	requestTimeout     = 10 * time.Second
)

type APIClient struct {
	logger  *zap.Logger
	client  *http.Client
	baseURL string
}

func NewClient(logger *zap.Logger) *APIClient {
	return &APIClient{
		logger: logger,
		client: &http.Client{
			Timeout: requestTimeout,
		},
		baseURL: baseURL,
	}
}

type chuckSearchResponse struct {
	Total  int `json:"total"`
	Result []struct {
		Categories []any  `json:"categories"`
		CreatedAt  string `json:"created_at"`
		IconURL    string `json:"icon_url"`
		ID         string `json:"id"`
		UpdatedAt  string `json:"updated_at"`
		URL        string `json:"url"`
		Value      string `json:"value"`
	} `json:"result"`
}

// Search makes a call to the chuck norris API search endpoint. The limit parameter is used
// to restrict what the caller gets -- the chuck norris endpoint does not support limits and
// can return large results (e.g. 9667 records for a query of 'chuck').
func (c *APIClient) Search(ctx context.Context, query string, limit int) ([]*domain.Joke, error) {
	logger := c.logger.With(zap.String("query", query), zap.Int("limit", limit))
	logger.Info("calling api search")

	url := fmt.Sprintf("%s/jokes/search?query=%s", c.baseURL, url.QueryEscape(query))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "https://github.com/davemolk/chuck")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response status: %s", resp.Status)
	}

	var data chuckSearchResponse
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode resp body: %w", err)
	}

	logger.Info("successful call", zap.Int("count", data.Total))

	if data.Total == 0 {
		logger.Debug("no results")
		return nil, nil
	}

	if limit > data.Total {
		limit = data.Total
		logger.Debug("requested more than we got, adjusting return", zap.Int("adjusted_limit", limit))
	}

	out := make([]*domain.Joke, limit)
	for i, j := range data.Result {
		if i == limit {
			break
		}

		// like its namesake, the chuck norris api is unpredictable, returning
		// varying decimal precision for the created_at field, so we will just
		// strip it out entirely
		if idx := strings.Index(j.CreatedAt, "."); idx != -1 {
			j.CreatedAt = j.CreatedAt[:idx]
		}

		createdAt, err := time.Parse(responseTimeFormat, j.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to parse time from api response: %w", err)
		}

		joke := &domain.Joke{
			ExternalID: j.ID,
			Content:    j.Value,
			URL:        j.URL,
			CreatedAt:  createdAt,
		}

		out[i] = joke
	}

	return out, nil
}

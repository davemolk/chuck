package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/davemolk/chuck/internal/domain"
	"github.com/davemolk/chuck/internal/tests/fixture"
	"github.com/davemolk/chuck/internal/tests/mock"
	"github.com/stretchr/testify/require"
)

func TestGetRandomJoke(t *testing.T) {
	var gotCtx context.Context
	jokeService := &mock.JokeService{
		GetRandomJokeFn: func(ctx context.Context) (*domain.Joke, error) {
			gotCtx = ctx
			return nil, errors.New("nope")
		},
	}

	h := NewJokeHandlers(fixture.TestLogger(t), jokeService)
	t.Run("handle service error", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/api/v1/jokes/random", nil)

		h.GetRandomJoke(w, r)
		require.True(t, jokeService.GetRandomJokeCalled)
		require.Equal(t, http.StatusInternalServerError, w.Code)
		jokeService.ResetCalls()
	})

	t.Run("success", func(t *testing.T) {
		jokeService = &mock.JokeService{
			GetRandomJokeFn: func(ctx context.Context) (*domain.Joke, error) {
				gotCtx = ctx
				return &domain.Joke{
					Content: "beard",
				}, nil
			},
		}
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/api/v1/jokes/random", nil)
		h.jokeService = jokeService

		h.GetRandomJoke(w, r)
		require.True(t, jokeService.GetRandomJokeCalled)
		require.NotNil(t, gotCtx)

		require.Equal(t, http.StatusOK, w.Code)

		var got domain.Joke
		err := json.NewDecoder(w.Body).Decode(&got)
		require.NoError(t, err)

		require.Equal(t, "beard", got.Content)
	})

}

func TestGetRandomJokeByQuery(t *testing.T) {
	var gotCtx context.Context
	var gotQuery string
	query := "beard"
	jokeService := &mock.JokeService{
		GetRandomJokeByQueryFn: func(ctx context.Context, query string) (*domain.Joke, error) {
			gotCtx = ctx
			gotQuery = query
			return nil, errors.New("blah")
		},
	}
	h := NewJokeHandlers(fixture.TestLogger(t), jokeService)

	t.Run("query required", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/api/v1/jokes/search?query=", nil)

		h.GetRandomJokeByQuery(w, r)

		require.False(t, jokeService.GetRandomJokeByQueryCalled)
		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("handle service error", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/jokes/search?query=%s", query), nil)

		h.GetRandomJokeByQuery(w, r)

		require.True(t, jokeService.GetRandomJokeByQueryCalled)
		require.Equal(t, http.StatusInternalServerError, w.Code)
		jokeService.ResetCalls()
	})

	t.Run("success", func(t *testing.T) {
		jokeService = &mock.JokeService{
			GetRandomJokeByQueryFn: func(ctx context.Context, query string) (*domain.Joke, error) {
				gotCtx = ctx
				gotQuery = query
				return &domain.Joke{
					Content: "beard",
				}, nil
			},
		}
		h.jokeService = jokeService

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/jokes/search?query=%s", query), nil)

		h.GetRandomJokeByQuery(w, r)

		require.True(t, jokeService.GetRandomJokeByQueryCalled)
		require.NotNil(t, gotCtx)

		require.Equal(t, http.StatusOK, w.Code)
		require.Equal(t, query, gotQuery)

		var got domain.Joke
		err := json.NewDecoder(w.Body).Decode(&got)
		require.NoError(t, err)

		require.Equal(t, "beard", got.Content)
	})
}

func TestGetPersonalizedJoke(t *testing.T) {
	var gotCtx context.Context
	var gotName string
	name := "dave"
	jokeService := &mock.JokeService{
		GetPersonalizedJokeFn: func(ctx context.Context, name string) (*domain.Joke, error) {
			gotCtx = ctx
			gotName = name
			return nil, errors.New("blah")
		},
	}
	h := NewJokeHandlers(fixture.TestLogger(t), jokeService)

	t.Run("name required", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/api/v1/jokes/personalized?name=", nil)

		h.GetPersonalizedJoke(w, r)

		require.False(t, jokeService.GetPersonalizedJokeCalled)
		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("handle service error", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/jokes/personalized?name=%s", name), nil)

		h.GetPersonalizedJoke(w, r)

		require.True(t, jokeService.GetPersonalizedJokeCalled)
		require.Equal(t, http.StatusInternalServerError, w.Code)
		jokeService.ResetCalls()
	})

	t.Run("success", func(t *testing.T) {
		jokeService = &mock.JokeService{
			GetPersonalizedJokeFn: func(ctx context.Context, name string) (*domain.Joke, error) {
				gotCtx = ctx
				gotName = name
				return &domain.Joke{
					Content: name,
				}, nil
			},
		}
		h.jokeService = jokeService

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/jokes/personalized?name=%s", name), nil)

		h.GetPersonalizedJoke(w, r)

		require.True(t, jokeService.GetPersonalizedJokeCalled)
		require.NotNil(t, gotCtx)

		require.Equal(t, http.StatusOK, w.Code)
		require.Equal(t, name, gotName)

		var got domain.Joke
		err := json.NewDecoder(w.Body).Decode(&got)
		require.NoError(t, err)

		require.Equal(t, name, got.Content)
	})
}

package joke

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/davemolk/chuck/internal/domain"
	"github.com/davemolk/chuck/internal/sql/dbtest"
	"github.com/davemolk/chuck/internal/tests/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestPersonalize(t *testing.T) {
	s := &Service{}

	tests := []struct {
		name     string
		joke     string
		person   string
		expected string
	}{
		{
			name:     "replaces standard capitalization",
			joke:     "Chuck Norris can divide by zero.",
			person:   "Bob",
			expected: "Bob can divide by zero.",
		},
		{
			name:     "replaces lowercase",
			joke:     "chuck norris writes code that optimizes itself.",
			person:   "Alice",
			expected: "alice writes code that optimizes itself.",
		},
		{
			name:     "replaces uppercase",
			joke:     "CHUCK NORRIS DOES NOT SLEEP.",
			person:   "Bob",
			expected: "BOB DOES NOT SLEEP.",
		},
		{
			name:     "replaces possessive with s",
			joke:     "Chuck Norris's keyboard has no escape key.",
			person:   "Alice",
			expected: "Alice's keyboard has no escape key.",
		},
		{
			name:   "replaces possessive with apostrophe only",
			joke:   "Chuck Norris' code never has bugs.",
			person: "Bob",
			// todo change this if time
			expected: "Bob' code never has bugs.",
		},
		{
			name:     "replaces multiple occurrences",
			joke:     "Chuck Norris met chuck norris and CHUCK NORRIS.",
			person:   "Alice",
			expected: "Alice met alice and ALICE.",
		},
		{
			name:     "no Chuck Norris present",
			joke:     "This joke has no reference.",
			person:   "Bob",
			expected: "This joke has no reference.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := s.personalize(tt.joke, tt.person)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestGetRandomDBJoke(t *testing.T) {
	db := dbtest.SetupTestDB(t)
	s := NewService(zap.NewNop(), db, nil)
	ctx := context.Background()

	joke, err := s.getRandomDBJoke(ctx)
	require.NoError(t, err)
	require.Contains(t, joke.Content, "Chuck")
}

func TestGetPersonalizedJoke(t *testing.T) {
	db := dbtest.SetupTestDB(t)
	s := NewService(zap.NewNop(), db, nil)
	ctx := context.Background()

	name := "dave molk"
	joke, err := s.GetPersonalizedJoke(ctx, name)
	require.NoError(t, err)
	require.Contains(t, joke.Content, name)
}

func TestGetRandomDBJokeByQuery(t *testing.T) {
	db := dbtest.SetupTestDB(t)
	s := NewService(zap.Must(zap.NewDevelopment()), db, nil)
	ctx := context.Background()

	t.Run("success, joke in db", func(t *testing.T) {
		joke, err := s.getRandomDBJokeByQuery(ctx, "uniq")
		require.NoError(t, err)
		require.Equal(t, int64(4), joke.ID)
	})

	t.Run("success, get from api", func(t *testing.T) {
		client := &mock.ChuckClient{
			SearchFn: func(ctx context.Context, query string, limit int) ([]*domain.Joke, error) {
				return []*domain.Joke{
					{
						ExternalID: "h0VbNVqJQcWQvpxtimWJ7Q",
						URL:        "https://api.chucknorris.io/jokes/h0VbNVqJQcWQvpxtimWJ7Q",
						Content:    "school didnt teach Chuck Norris he taught school",
						CreatedAt:  time.Now(),
					},
				}, nil
			},
		}
		s.client = client
		joke, err := s.GetRandomJokeByQuery(ctx, "school")
		require.NoError(t, err)
		require.Equal(t, int64(5), joke.ID)
	})

	t.Run("success, get from api with multiple inserts", func(t *testing.T) {
		client := &mock.ChuckClient{
			SearchFn: func(ctx context.Context, query string, limit int) ([]*domain.Joke, error) {
				return []*domain.Joke{
					// we will include the same joke as previous test to make sure we handle
					// insert conflicts correctly
					{
						ExternalID: "h0VbNVqJQcWQvpxtimWJ7Q",
						URL:        "https://api.chucknorris.io/jokes/h0VbNVqJQcWQvpxtimWJ7Q",
						Content:    "school didnt teach Chuck Norris he taught school",
						CreatedAt:  time.Now(),
					},
					{
						ExternalID: "z_VZSvW5SWud7-Vb0oZgIw",
						URL:        "https://api.chucknorris.io/jokes/z_VZSvW5SWud7-Vb0oZgIw",
						Content:    "The leading cause of ninja death is Chuck Norris.",
						CreatedAt:  time.Now(),
					},
					{
						ExternalID: "Yf3aq8BRSQmL9NBwWPoFqA",
						URL:        "https://api.chucknorris.io/jokes/Yf3aq8BRSQmL9NBwWPoFqA",
						Content:    "Chuck Norris is the reason ninjas like to be unseen.",
						CreatedAt:  time.Now(),
					},
					{
						ExternalID: "wIkJ7EssS6GUs-ACNUbzuw",
						URL:        "https://api.chucknorris.io/jokes/wIkJ7EssS6GUs-ACNUbzuw",
						Content:    "There is no Ninja Turtle cereal because eating ninjas for breakfast is a copyright of Chuck Norris.",
						CreatedAt:  time.Now(),
					},
				}, nil
			},
		}
		s.client = client
		joke, err := s.GetRandomJokeByQuery(ctx, "ninja")
		require.NoError(t, err)
		externalIDs := []string{"h0VbNVqJQcWQvpxtimWJ7Q", "z_VZSvW5SWud7-Vb0oZgIw", "Yf3aq8BRSQmL9NBwWPoFqA", "wIkJ7EssS6GUs-ACNUbzuw"}
		// make sure joke is from one of the ones returned by api
		require.Contains(t, externalIDs, joke.ExternalID)
	})

	t.Run("returns ErrNoJokes when api returns no results", func(t *testing.T) {
		client := &mock.ChuckClient{
			SearchFn: func(ctx context.Context, query string, limit int) ([]*domain.Joke, error) {
				return []*domain.Joke{}, nil
			},
		}
		s.client = client
		_, err := s.GetRandomJokeByQuery(ctx, "kale")
		require.Error(t, err)
		require.True(t, errors.Is(err, ErrNoJokes))
	})
}

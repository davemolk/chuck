package joke

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/rand/v2"
	"strings"

	"github.com/davemolk/chuck/internal/domain"
	"github.com/davemolk/chuck/internal/service"
	sqldb "github.com/davemolk/chuck/internal/sql"
	"go.uber.org/zap"
)

const maxJokesFromAPI = 100

var _ service.JokeService = (*Service)(nil)

var ErrNoJokes = errors.New("no jokes found")

type chuckGetter interface {
	Search(ctx context.Context, query string, limit int) ([]*domain.Joke, error)
}

type Service struct {
	logger *zap.Logger
	db     *sqldb.DB
	client chuckGetter
}

func NewService(logger *zap.Logger, db *sqldb.DB, client chuckGetter) *Service {
	return &Service{
		logger: logger,
		db:     db,
		client: client,
	}
}

func (s *Service) GetPersonalizedJoke(ctx context.Context, name string) (*domain.Joke, error) {
	joke, err := s.GetRandomJoke(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get joke for personalization: %w", err)
	}

	personalized := s.personalize(joke.Content, name)
	joke.Content = personalized

	return joke, nil
}

// GetRandomJoke selects a joke at random from the database. Since we seed in
// the initial migration, we will always have a result.
func (s *Service) GetRandomJoke(ctx context.Context) (*domain.Joke, error) {
	query := `
		select id, external_id, joke_url, content, created_at
		from jokes
		order by random()
		limit 1
	`

	var joke domain.Joke
	err := s.db.QueryRowContext(ctx, query).Scan(
		&joke.ID, &joke.ExternalID, &joke.URL, &joke.Content, &joke.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get joke: %w", err)
	}

	return &joke, nil
}

// personalize is a quick and dirty (but readable) substitution of
// a name for some variant of Chuck Norris.
func (s *Service) personalize(joke, name string) string {
	joke = strings.ReplaceAll(joke, "Chuck Norris", name)
	joke = strings.ReplaceAll(joke, "chuck norris", strings.ToLower(name))
	joke = strings.ReplaceAll(joke, "CHUCK NORRIS", strings.ToUpper(name))

	// handle possessive where name doesn't end with s, e.g. we want Bob's, not Bob'
	if !strings.HasSuffix(name, "s") && strings.Contains(joke, name+"' ") {
		joke = strings.ReplaceAll(joke, name+"' ", name+"'s ")
	}

	return joke
}

// GetRandomJokeByQuery searches the database for a joke whose content matches the
// query. If this fails, GetRandomJokeByQuery calls the search endpoint of the Chuck
// Norris API, saving any results to the database and returning one to user.
func (s *Service) GetRandomJokeByQuery(ctx context.Context, query string) (*domain.Joke, error) {
	// first, check database for a match
	joke, err := s.getRandomDBJokeByQuery(ctx, query)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return nil, err
	}

	// exit early if we found our match
	if joke != nil {
		return joke, nil
	}

	// if we didn't match, call the api
	logger := s.logger.With(zap.String("query", query))
	logger.Info("no cached matches, calling api...")

	jokes, err := s.client.Search(ctx, query, maxJokesFromAPI)
	if err != nil {
		return nil, fmt.Errorf("failed to search api: %w", err)
	}

	if len(jokes) == 0 {
		return nil, ErrNoJokes
	}

	logger.Info("found", zap.Int("count", len(jokes)))

	// populate db. as we insert, we will scan the id to the slice we're passing in
	if err = s.saveJokes(ctx, jokes); err != nil {
		// user should still get a joke if we can't save them
		logger.Error("failed to save jokes", zap.Error(err))
	}

	// jokes will now have the id, so we can return from the slice instead
	// of hitting the db again
	return jokes[rand.IntN(len(jokes))], nil
}

func (s *Service) getRandomDBJokeByQuery(ctx context.Context, query string) (*domain.Joke, error) {
	q := `
		select id, external_id, joke_url, content, created_at
		from jokes
		where to_tsvector('simple', content) @@ to_tsquery('simple', $1)
		order by random()
		limit 1
	`

	var joke domain.Joke
	err := s.db.QueryRowContext(ctx, q, query).Scan(
		&joke.ID, &joke.ExternalID, &joke.URL, &joke.Content, &joke.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get joke by content: %w", err)
	}

	return &joke, nil
}

func (s *Service) saveJokes(ctx context.Context, jokes []*domain.Joke) error {
	// sanity check
	if len(jokes) == 0 {
		return nil
	}

	// this presupposes that external id is a unique identifier, and it certainly
	// appears to be, but there's no statement to that effect in the chuck norris
	// api docs.
	//
	// since we need to scan an id and we might have conflicts, we could either check
	// for sql.ErrNoRows and handle, or do a no-op update which will return the id.
	query := `
		insert into jokes (external_id, joke_url, content, created_at)
		values ($1, $2, $3, $4)
		on conflict (external_id) do update
		set external_id = jokes.external_id
		returning id
	`

	return s.db.RunInTx(ctx, func(tx *sql.Tx) error {
		for _, joke := range jokes {
			if err := tx.QueryRowContext(ctx, query, joke.ExternalID, joke.URL, joke.Content, joke.CreatedAt).Scan(&joke.ID); err != nil {
				return err
			}
		}

		return nil
	})
}

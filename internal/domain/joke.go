package domain

import (
	"errors"
	"time"
)

var ErrNotFound = errors.New("record not found")

type Joke struct {
	ID         int64
	ExternalID string
	Content    string
	URL        string
	CreatedAt  time.Time
}

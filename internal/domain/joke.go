package domain

import (
	"time"

	"github.com/gofrs/uuid"
)

type Joke struct {
	ID         uuid.UUID
	ExternalID string
	Content    string
	URL        string
	CreatedAt  time.Time
}

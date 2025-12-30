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

type Token struct {
	Plaintext string `json:"token"`
	Hash      []byte `json:"-"`
	UserID    int64  `json:"-"`
	ExpiresAt time.Time
}

type User struct {
	ID        int64     `json:"id"`
	HashedPW  []byte    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	Email     string    `json:"email"`
}

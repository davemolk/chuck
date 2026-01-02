package domain

import (
	"errors"
	"time"
)

var ErrNotFound = errors.New("record not found")

type Joke struct {
	ID         int64     `json:"id"`
	ExternalID string    `json:"external_id"`
	Content    string    `json:"content"`
	URL        string    `json:"original_url"`
	CreatedAt  time.Time `json:"creatd_at"`
}

type Token struct {
	Plaintext string    `json:"token"`
	Hash      []byte    `json:"-"`
	UserID    int64     `json:"-"`
	ExpiresAt time.Time `json:"expires_at"`
}

type User struct {
	ID        int64     `json:"id"`
	HashedPW  []byte    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	Email     string    `json:"email"`
}

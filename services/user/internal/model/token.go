package model

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"time"
)

type Scope string

const (
	ScopeAuthentication Scope = "authentication"
)

const tokenLength = 32

type Token struct {
	Bytes     []byte
	Text      string
	UserID    int64
	Scope     Scope
	ExpiresAt time.Time
}

func NewToken(userID int64, scope Scope, ttl time.Duration) (*Token, error) {
	bytes := make([]byte, tokenLength)

	_, err := rand.Read(bytes)
	if err != nil {
		return nil, err
	}

	return &Token{
		Bytes:     bytes,
		Text:      base64.URLEncoding.EncodeToString(bytes),
		UserID:    userID,
		ExpiresAt: time.Now().Add(ttl),
		Scope:     scope,
	}, nil
}

func (t *Token) DecodeBytes() error {
	bytes, err := base64.URLEncoding.DecodeString(t.Text)
	if err != nil {
		return err
	}
	t.Bytes = bytes
	return nil
}

type TokenModel struct {
	DB *sql.DB
}

func (tm TokenModel) Insert(token *Token) error {
	query := `
	INSERT INTO tokens (hash, user_id, scope, expires_at)
	VALUES ($1, $2, $3, $4);`

	args := []any{token.Bytes, token.UserID, token.Scope, token.ExpiresAt}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := tm.DB.ExecContext(ctx, query, args...)
	return err
}

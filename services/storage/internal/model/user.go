package model

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
)

var ErrDuplicateEmail = errors.New("user with this email address already exists")

type User struct {
	ID           int64
	Name         string
	Email        string
	PasswordHash []byte
	RoleID       int64
	CreatedAt    time.Time
	Version      int
}

type UserModel struct {
	DB *sql.DB
}

func (um UserModel) Insert(user *User) error {
	query := `
	INSERT INTO USERS (name, email, password_hash, created_at, role_id)
	VALUES ($1, $2, $3, $4, (SELECT id FROM roles WHERE name='User'))
	RETURNING id`
	args := []any{user.Name, user.Email, user.PasswordHash, user.CreatedAt}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := um.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code {
			case "23505":
				return ErrDuplicateEmail
			default:
				return err
			}
		}
	}

	return nil
}

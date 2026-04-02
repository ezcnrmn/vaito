package model

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail = errors.New("user with this email address already exists")
	ErrUserNotFound   = errors.New("user not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Role struct {
	Name string
	ID   int64
}
type Password struct {
	Hash      []byte
	Plaintext string
}

func (p *Password) GenerateHash() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(p.Plaintext), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	p.Hash = hash
	return nil
}

func (p Password) CompareHash() bool {
	return nil == bcrypt.CompareHashAndPassword(p.Hash, []byte(p.Plaintext))
}

type User struct {
	ID        int64
	Name      string
	Email     string
	Password  Password
	Role      Role
	CreatedAt time.Time
	Version   int
}

type UserModel struct {
	DB *sql.DB
}

func (um UserModel) Insert(user *User) error {
	query := `
	INSERT INTO users (name, email, password_hash, created_at, role_id)
	VALUES ($1, $2, $3, $4, (SELECT id FROM roles WHERE name='User'))
	RETURNING id, (SELECT name FROM roles WHERE id=role_id) AS role_name;`

	args := []any{user.Name, user.Email, user.Password.Hash, user.CreatedAt}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := um.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.Role.Name)
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

func (um UserModel) GetUser(id int64) (*User, error) {
	query := `
	SELECT u.id, u.name, u.email, r.name, u.version
	FROM users u
	JOIN roles r ON u.role_id = r.id
	WHERE u.id = $1;`

	args := []any{id}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	user := &User{}
	err := um.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.Name, &user.Email, &user.Role.Name, &user.Version)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrUserNotFound
		default:
			return nil, err
		}
	}

	return user, nil
}

func (um UserModel) GetUserByEmail(email string) (*User, error) {
	query := `
	SELECT id, name, email, role_id, password_hash, version
	FROM users
	WHERE email = $1;`

	args := []any{email}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	user := &User{}
	err := um.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.Name, &user.Email, &user.Role.ID, &user.Password.Hash, &user.Version)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrUserNotFound
		default:
			return nil, err
		}
	}

	return user, nil
}

func (um UserModel) GetUserByToken(token *Token) (*User, error) {
	query := `
	SELECT u.id, u.name, u.email, r.name, u.password_hash, u.version
	FROM users u
	JOIN roles r ON u.role_id = r.id
	JOIN tokens t ON u.id = t.user_id
	WHERE t.hash = $1
		AND t.scope = $2
		AND t.expires_at > $3;`

	args := []any{token.Bytes, token.Scope, time.Now()}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	user := &User{}
	err := um.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.Name, &user.Email, &user.Role.Name, &user.Password.Hash, &user.Version)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrUserNotFound
		default:
			return nil, err
		}
	}

	return user, nil
}

func (um UserModel) GetUserIDByToken(token *Token) (int64, error) {
	query := `
	SELECT user_id
	FROM tokens
	WHERE hash = $1 AND scope = $2 AND expires_at > $3`

	args := []any{token.Bytes, token.Scope, time.Now()}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var userID int64
	err := um.DB.QueryRowContext(ctx, query, args...).Scan(&userID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return 0, ErrUserNotFound
		default:
			return 0, err
		}
	}

	return userID, nil
}

func (um UserModel) GetUserPermissionsByToken(token *Token) (int64, []string, error) {
	query := `
	SELECT DISTINCT u.id, p.code
	FROM tokens t
	JOIN users u ON u.id = t.user_id
	LEFT JOIN roles r ON u.role_id = r.id
	LEFT JOIN roles_permissions rp ON r.id = rp.role_id
	LEFT JOIN permissions p ON rp.permission_id = p.id
	WHERE t.hash = $1 AND t.scope = $2 AND t.expires_at > $3`

	args := []any{token.Bytes, token.Scope, time.Now()}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := um.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return 0, nil, err
	}
	defer rows.Close()

	var userID int64
	var permissions []string
	var found bool
	for rows.Next() {
		found = true
		var p sql.NullString
		err := rows.Scan(&userID, &p)
		if err != nil {
			return 0, nil, err
		}
		permissions = append(permissions, p.String)
	}
	if err = rows.Err(); err != nil {
		return 0, nil, err
	}

	if !found {
		return 0, nil, ErrUserNotFound
	}

	return userID, permissions, nil
}

func (um UserModel) UpdateUser(user *User) error {
	query := `
	UPDATE users
	SET name = $1, email = $2, password_hash = $3, version = version+1
	WHERE id = $4 AND version = $5
	RETURNING version`

	args := []any{user.Name, user.Email, user.Password.Hash, user.ID, user.Version}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := um.DB.QueryRowContext(ctx, query, args...).Scan(&user.Version)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

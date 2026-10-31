package store

import (
	"context"
	"database/sql"
	"errors"
)

// User is an account. The hash never leaves the server.
type User struct {
	ID           int64  `json:"id"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
}

// CreateUser inserts a user with an already-hashed password.
func (s *Store) CreateUser(
	ctx context.Context, email, hash string,
) (User, error) {
	res, err := s.db.ExecContext(ctx,
		`INSERT INTO users (email, password_hash)
		 VALUES (?, ?)`,
		email, hash)
	if isUnique(err) {
		return User{}, ErrDuplicate
	}
	if err != nil {
		return User{}, err
	}
	id, err := res.LastInsertId()
	return User{ID: id, Email: email, PasswordHash: hash}, err
}

// UserByEmail finds a user for login.
func (s *Store) UserByEmail(
	ctx context.Context, email string,
) (User, error) {
	var u User
	err := s.db.QueryRowContext(ctx,
		`SELECT id, email, password_hash
		 FROM users WHERE email = ?`, email).
		Scan(&u.ID, &u.Email, &u.PasswordHash)
	if errors.Is(err, sql.ErrNoRows) {
		return u, ErrNotFound
	}
	return u, err
}

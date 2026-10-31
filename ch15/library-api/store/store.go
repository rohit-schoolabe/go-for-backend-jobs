// Package store holds all the SQL for the library API.
package store

import (
	"database/sql"
	_ "embed"
	"errors"
	"strings"

	_ "modernc.org/sqlite"
)

// Errors the API layer turns into status codes.
var (
	ErrNotFound  = errors.New("not found")
	ErrDuplicate = errors.New("already exists")
	ErrOnLoan    = errors.New("book is already on loan")
)

//go:embed schema.sql
var schema string

// Store wraps the database connection pool.
type Store struct {
	db *sql.DB
}

// Open connects to a SQLite file, or ":memory:",
// and creates the tables.
func Open(dsn string) (*Store, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	// SQLite has one writer. One connection also keeps
	// a ":memory:" database shared by every query.
	db.SetMaxOpenConns(1)

	_, err = db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(schema); err != nil {
		return nil, err
	}
	return &Store{db: db}, nil
}

// Close releases the database.
func (s *Store) Close() error {
	return s.db.Close()
}

func isUnique(err error) bool {
	return err != nil &&
		strings.Contains(err.Error(), "UNIQUE")
}

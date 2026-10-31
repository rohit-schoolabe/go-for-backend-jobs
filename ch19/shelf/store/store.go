// Package store holds all the SQL for the library API.
// It runs on SQLite for local work and on PostgreSQL in
// production, chosen by the DSN.
package store

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib" // driver "pgx"
	_ "modernc.org/sqlite"             // driver "sqlite"
)

// Errors the API layer turns into status codes.
var (
	ErrNotFound  = errors.New("not found")
	ErrDuplicate = errors.New("already exists")
	ErrOnLoan    = errors.New("book is already on loan")
)

// Store wraps the database connection pool.
type Store struct {
	db *sql.DB
	pg bool // PostgreSQL rather than SQLite
}

// Open connects to the database named by dsn and applies
// any migrations it has not seen. A postgres:// URL picks
// PostgreSQL; anything else is a SQLite file or ":memory:".
func Open(dsn string) (*Store, error) {
	pg := strings.HasPrefix(dsn, "postgres://") ||
		strings.HasPrefix(dsn, "postgresql://")
	driver := "sqlite"
	if pg {
		driver = "pgx"
	}
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}
	s := &Store{db: db, pg: pg}
	if pg {
		db.SetMaxOpenConns(10)
		db.SetConnMaxIdleTime(5 * time.Minute)
	} else {
		db.SetMaxOpenConns(1)
		_, err = db.Exec("PRAGMA foreign_keys = ON")
		if err != nil {
			db.Close()
			return nil, err
		}
	}
	ctx, cancel := context.WithTimeout(
		context.Background(), 30*time.Second)
	defer cancel()
	if err := s.migrate(ctx); err != nil {
		db.Close()
		return nil, err
	}
	return s, nil
}

// Close releases the database.
func (s *Store) Close() error {
	return s.db.Close()
}

// Ping reports whether the database answers. The
// readiness check calls it.
func (s *Store) Ping(ctx context.Context) error {
	return s.db.PingContext(ctx)
}

// q rewrites ? placeholders as $1, $2 ... for PostgreSQL.
// The queries in this package never put ? inside a string.
func (s *Store) q(query string) string {
	if !s.pg {
		return query
	}
	var b strings.Builder
	n := 0
	for _, r := range query {
		if r == '?' {
			n++
			b.WriteString("$" + strconv.Itoa(n))
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}

// isUnique spots a unique-constraint failure from
// either driver.
func isUnique(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505" // unique_violation
	}
	return err != nil &&
		strings.Contains(err.Error(), "UNIQUE")
}

package store

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"path"
	"strings"
)

//go:embed migrations
var migrations embed.FS

// migrate applies, in name order, every file under
// migrations/<dialect> that schema_migrations does not
// list yet. All of it happens in one transaction.
func (s *Store) migrate(ctx context.Context) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if s.pg {
		// Two replicas starting together wait here, one
		// after the other, instead of racing.
		_, err := tx.ExecContext(ctx,
			"SELECT pg_advisory_xact_lock(5150)")
		if err != nil {
			return err
		}
	}
	_, err = tx.ExecContext(ctx,
		`CREATE TABLE IF NOT EXISTS schema_migrations (
		   version TEXT PRIMARY KEY)`)
	if err != nil {
		return err
	}

	dir := "migrations/sqlite"
	if s.pg {
		dir = "migrations/postgres"
	}
	files, err := fs.ReadDir(migrations, dir)
	if err != nil {
		return err
	}
	for _, f := range files { // ReadDir sorts by name
		version := strings.TrimSuffix(f.Name(), ".sql")
		var seen int
		err := tx.QueryRowContext(ctx, s.q(
			`SELECT COUNT(*) FROM schema_migrations
			 WHERE version = ?`), version).Scan(&seen)
		if err != nil {
			return err
		}
		if seen > 0 {
			continue
		}
		name := path.Join(dir, f.Name())
		body, err := migrations.ReadFile(name)
		if err != nil {
			return err
		}
		_, err = tx.ExecContext(ctx, string(body))
		if err != nil {
			return fmt.Errorf("%s: %w", name, err)
		}
		_, err = tx.ExecContext(ctx, s.q(
			`INSERT INTO schema_migrations (version)
			 VALUES (?)`), version)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

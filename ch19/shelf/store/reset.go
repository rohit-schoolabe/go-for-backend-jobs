package store

import "context"

// ResetForTest empties every table and restarts the IDs.
// Tests that share one PostgreSQL database call it first.
func (s *Store) ResetForTest(ctx context.Context) error {
	if !s.pg {
		return nil // each SQLite test opens ":memory:"
	}
	_, err := s.db.ExecContext(ctx,
		`TRUNCATE loans, books, users RESTART IDENTITY`)
	return err
}

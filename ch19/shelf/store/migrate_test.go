package store

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

// Opening the same database twice must not apply a
// migration twice; every container restart does this.
func TestMigrateTwice(t *testing.T) {
	dsn := os.Getenv("SHELF_TEST_DSN")
	if dsn == "" {
		dsn = filepath.Join(t.TempDir(), "shelf.db")
	}
	for i := 0; i < 2; i++ {
		st, err := Open(dsn)
		if err != nil {
			t.Fatalf("open %d: %v", i+1, err)
		}
		var n int
		err = st.db.QueryRowContext(context.Background(),
			`SELECT COUNT(*) FROM schema_migrations`).Scan(&n)
		st.Close()
		if err != nil || n != 1 {
			t.Fatalf("open %d: %d rows, err %v", i+1, n, err)
		}
	}
}

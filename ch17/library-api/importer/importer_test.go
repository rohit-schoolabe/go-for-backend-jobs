package importer

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/rohit-schoolabe/go-for-backend-jobs/ch16/library-api/store"
)

func export(n int) string {
	var sb strings.Builder
	sb.WriteString("[")
	for i := range n {
		if i > 0 {
			sb.WriteString(",")
		}
		year := 2000 + i%20
		if i%10 == 0 {
			year = 0 // the old system's missing years
		}
		fmt.Fprintf(&sb, `{"code":"B%d","title":" Book %d ",`+
			`"author":"Rao","year":%d}`, i, i, year)
	}
	sb.WriteString("]")
	return sb.String()
}

func TestImportIntoStore(t *testing.T) {
	st, err := store.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	ctx := context.Background()
	u, err := st.CreateUser(ctx, "admin@college.in", "x")
	if err != nil {
		t.Fatal(err)
	}

	stats, err := Import(ctx, strings.NewReader(export(200)),
		st, u.ID, 4)
	if err != nil {
		t.Fatal(err)
	}
	if stats.Imported != 180 || stats.Skipped != 20 {
		t.Fatalf("got %+v", stats)
	}
	// Fan-out does not keep the file's order, so check
	// membership, not position.
	books, _ := st.Books(ctx, "")
	titles := map[string]bool{}
	for _, b := range books {
		titles[b.Title] = true
	}
	if len(books) != 180 || !titles["Book 1"] {
		t.Fatalf("%d books, Book 1 present: %v",
			len(books), titles["Book 1"])
	}
}

type failingDB struct{ after int }

func (f *failingDB) CreateBook(context.Context,
	store.Book) (store.Book, error) {
	if f.after == 0 {
		return store.Book{}, errors.New("disk full")
	}
	f.after--
	return store.Book{}, nil
}

func TestImportStopsAllStagesOnError(t *testing.T) {
	before := runtime.NumGoroutine()
	_, err := Import(context.Background(),
		strings.NewReader(export(5000)),
		&failingDB{after: 3}, 1, 8)
	if err == nil || err.Error() != "disk full" {
		t.Fatalf("err = %v", err)
	}
	// Every stage must have exited: no leaked goroutines.
	time.Sleep(50 * time.Millisecond)
	if after := runtime.NumGoroutine(); after > before {
		t.Errorf("goroutines: %d before, %d after",
			before, after)
	}
}

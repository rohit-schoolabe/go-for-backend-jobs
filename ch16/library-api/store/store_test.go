package store

import (
	"context"
	"errors"
	"testing"
	"time"
)

// newStore gives each test its own empty database.
func newStore(t *testing.T) *Store {
	t.Helper()
	st, err := Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { st.Close() })
	return st
}

// seed adds one user and one book and returns their IDs.
func seed(t *testing.T, st *Store) (userID, bookID int64) {
	t.Helper()
	ctx := context.Background()
	u, err := st.CreateUser(ctx, "asha@college.in", "x")
	if err != nil {
		t.Fatal(err)
	}
	b, err := st.CreateBook(ctx, Book{Title: "Optics",
		Author: "Iyer", Year: 2019, AddedBy: u.ID})
	if err != nil {
		t.Fatal(err)
	}
	return u.ID, b.ID
}

func TestBorrowTwice(t *testing.T) {
	st := newStore(t)
	ctx := context.Background()
	uid, bid := seed(t, st)
	now := time.Now().UTC()

	if _, err := st.Borrow(ctx, bid, uid, now, 14); err != nil {
		t.Fatal(err)
	}
	_, err := st.Borrow(ctx, bid, uid, now, 14)
	if !errors.Is(err, ErrOnLoan) {
		t.Fatalf("second borrow: got %v, want ErrOnLoan", err)
	}
}

func TestOverdueAt(t *testing.T) {
	st := newStore(t)
	ctx := context.Background()
	uid, bid := seed(t, st)
	issued := time.Date(2026, 1, 1, 9, 0, 0, 0, time.UTC)
	_, err := st.Borrow(ctx, bid, uid, issued, 14)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("not yet due", func(t *testing.T) {
		list, err := st.OverdueAt(ctx, issued.AddDate(0, 0, 10))
		if err != nil || len(list) != 0 {
			t.Fatalf("got %v, %v", list, err)
		}
	})
	t.Run("overdue", func(t *testing.T) {
		list, err := st.OverdueAt(ctx, issued.AddDate(0, 0, 20))
		if err != nil || len(list) != 1 {
			t.Fatalf("got %v, %v", list, err)
		}
		if list[0].Email != "asha@college.in" {
			t.Errorf("email %q", list[0].Email)
		}
	})
	t.Run("returned", func(t *testing.T) {
		late := issued.AddDate(0, 0, 20)
		if err := st.Return(ctx, bid, uid, late); err != nil {
			t.Fatal(err)
		}
		list, _ := st.OverdueAt(ctx, late)
		if len(list) != 0 {
			t.Fatalf("returned loan still listed: %v", list)
		}
	})
}

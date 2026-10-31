package store

import (
	"context"
	"database/sql"
	"time"
)

// Loan is one borrowing of one book by one user.
type Loan struct {
	ID         int64      `json:"id"`
	BookID     int64      `json:"book_id"`
	BookTitle  string     `json:"book_title"`
	BorrowedAt time.Time  `json:"borrowed_at"`
	DueAt      time.Time  `json:"due_at"`
	ReturnedAt *time.Time `json:"returned_at,omitempty"`
}

// Borrow opens a loan. The one_open_loan index makes a
// second open loan on the same book fail.
func (s *Store) Borrow(
	ctx context.Context, bookID, userID int64,
	now time.Time, days int,
) (Loan, error) {
	due := now.AddDate(0, 0, days)
	loan := Loan{BookID: bookID, BorrowedAt: now,
		DueAt: due}
	err := s.db.QueryRowContext(ctx, s.q(
		`INSERT INTO loans
		   (book_id, user_id, borrowed_at, due_at)
		 VALUES (?, ?, ?, ?) RETURNING id`),
		bookID, userID, now, due).Scan(&loan.ID)
	if isUnique(err) {
		return Loan{}, ErrOnLoan
	}
	return loan, err
}

// Return closes the user's open loan on the book.
func (s *Store) Return(
	ctx context.Context, bookID, userID int64, now time.Time,
) error {
	res, err := s.db.ExecContext(ctx, s.q(
		`UPDATE loans SET returned_at = ?
		 WHERE book_id = ? AND user_id = ?
		   AND returned_at IS NULL`),
		now, bookID, userID)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return ErrNotFound
	}
	return nil
}

// LoansFor lists a user's loans, newest first.
func (s *Store) LoansFor(
	ctx context.Context, userID int64,
) ([]Loan, error) {
	rows, err := s.db.QueryContext(ctx, s.q(
		`SELECT l.id, l.book_id, b.title, l.borrowed_at,
		        l.due_at, l.returned_at
		 FROM loans l JOIN books b ON b.id = l.book_id
		 WHERE l.user_id = ? ORDER BY l.id DESC`), userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	loans := []Loan{}
	for rows.Next() {
		var l Loan
		var returned sql.NullTime
		err := rows.Scan(&l.ID, &l.BookID, &l.BookTitle,
			&l.BorrowedAt, &l.DueAt, &returned)
		if err != nil {
			return nil, err
		}
		if returned.Valid {
			l.ReturnedAt = &returned.Time
		}
		loans = append(loans, l)
	}
	return loans, rows.Err()
}

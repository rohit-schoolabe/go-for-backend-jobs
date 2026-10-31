package store

import (
	"context"
	"time"
)

// Overdue is an open loan past its due date, with what a
// reminder needs.
type Overdue struct {
	LoanID int64
	Email  string
	Title  string
	DueAt  time.Time
}

// OverdueAt lists open loans that were due before now.
func (s *Store) OverdueAt(
	ctx context.Context, now time.Time,
) ([]Overdue, error) {
	rows, err := s.db.QueryContext(ctx, s.q(
		`SELECT l.id, u.email, b.title, l.due_at
		 FROM loans l
		 JOIN users u ON u.id = l.user_id
		 JOIN books b ON b.id = l.book_id
		 WHERE l.returned_at IS NULL AND l.due_at < ?
		 ORDER BY l.due_at`), now)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Overdue
	for rows.Next() {
		var o Overdue
		err := rows.Scan(&o.LoanID, &o.Email, &o.Title,
			&o.DueAt)
		if err != nil {
			return nil, err
		}
		list = append(list, o)
	}
	return list, rows.Err()
}

package store

import (
	"context"
	"database/sql"
	"errors"
	"strings"
)

// Book is one title, with whether it is out right now.
type Book struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	Author  string `json:"author"`
	Year    int    `json:"year"`
	AddedBy int64  `json:"added_by"`
	OnLoan  bool   `json:"on_loan"`
}

const bookCols = `b.id, b.title, b.author, b.year, b.added_by,
  EXISTS (SELECT 1 FROM loans l
          WHERE l.book_id = b.id AND l.returned_at IS NULL)`

// scanner is satisfied by both *sql.Row and *sql.Rows.
type scanner interface {
	Scan(dest ...any) error
}

func scanBook(sc scanner) (Book, error) {
	var b Book
	err := sc.Scan(&b.ID, &b.Title, &b.Author, &b.Year,
		&b.AddedBy, &b.OnLoan)
	return b, err
}

// CreateBook inserts b and returns it with its new ID.
func (s *Store) CreateBook(
	ctx context.Context, b Book,
) (Book, error) {
	err := s.db.QueryRowContext(ctx, s.q(
		`INSERT INTO books (title, author, year, added_by)
		 VALUES (?, ?, ?, ?) RETURNING id`),
		b.Title, b.Author, b.Year, b.AddedBy).Scan(&b.ID)
	return b, err
}

// Books lists books, filtered by author when author is set.
func (s *Store) Books(
	ctx context.Context, author string,
) ([]Book, error) {
	q := `SELECT ` + bookCols + ` FROM books b`
	var args []any
	if author != "" {
		q += ` WHERE LOWER(b.author) LIKE ?`
		args = append(args, "%"+strings.ToLower(author)+"%")
	}
	q += ` ORDER BY b.id`
	rows, err := s.db.QueryContext(ctx, s.q(q), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	books := []Book{}
	for rows.Next() {
		b, err := scanBook(rows)
		if err != nil {
			return nil, err
		}
		books = append(books, b)
	}
	return books, rows.Err()
}

// Book fetches one book by ID.
func (s *Store) Book(
	ctx context.Context, id int64,
) (Book, error) {
	row := s.db.QueryRowContext(ctx, s.q(
		`SELECT `+bookCols+` FROM books b WHERE b.id = ?`), id)
	b, err := scanBook(row)
	if errors.Is(err, sql.ErrNoRows) {
		return b, ErrNotFound
	}
	return b, err
}

// UpdateBook saves new title, author and year.
func (s *Store) UpdateBook(ctx context.Context, b Book) error {
	_, err := s.db.ExecContext(ctx, s.q(
		`UPDATE books SET title = ?, author = ?, year = ?
		 WHERE id = ?`),
		b.Title, b.Author, b.Year, b.ID)
	return err
}

// DeleteBook removes a book and, by cascade, its loans.
func (s *Store) DeleteBook(
	ctx context.Context, id int64,
) error {
	_, err := s.db.ExecContext(ctx, s.q(
		`DELETE FROM books WHERE id = ?`), id)
	return err
}

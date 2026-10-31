// Package loan models one book lent to one member.
package loan

import (
	"errors"
	"fmt"
	"time"
)

// FinePerDay is charged for each day past the due date.
const FinePerDay = 2.0

// Loan is a book issued to a member. Its fields are unexported,
// so other packages must go through New and the methods.
type Loan struct {
	bookCode string
	member   string
	issuedAt time.Time
	dueAt    time.Time
}

// New validates its input and returns a loan due in 14 days.
func New(bookCode, member string,
	issuedAt time.Time) (*Loan, error) {
	if bookCode == "" || member == "" {
		return nil, errors.New("missing book code or member")
	}
	return &Loan{
		bookCode: bookCode,
		member:   member,
		issuedAt: issuedAt,
		dueAt:    issuedAt.AddDate(0, 0, 14),
	}, nil
}

// Describe returns a one-line summary. It does not change l.
func (l Loan) Describe() string {
	return fmt.Sprintf("%s lent to %s, due %s",
		l.bookCode, l.member, l.dueAt.Format("02 Jan 2006"))
}

// Renew moves the due date forward in place.
func (l *Loan) Renew(days int) {
	l.dueAt = l.dueAt.AddDate(0, 0, days)
}

// Fine is what the member owes if the book comes back at now.
func (l Loan) Fine(now time.Time) float64 {
	late := int(now.Sub(l.dueAt).Hours() / 24)
	if late <= 0 {
		return 0
	}
	return float64(late) * FinePerDay
}

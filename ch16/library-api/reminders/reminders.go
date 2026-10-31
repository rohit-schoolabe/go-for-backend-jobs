// Package reminders emails members whose books are late.
package reminders

import (
	"context"
	"fmt"
	"time"

	"github.com/rohit-schoolabe/go-for-backend-jobs/ch16/library-api/fines"
	"github.com/rohit-schoolabe/go-for-backend-jobs/ch16/library-api/store"
)

// Source lists overdue loans. *store.Store satisfies it.
type Source interface {
	OverdueAt(ctx context.Context,
		now time.Time) ([]store.Overdue, error)
}

// Sender delivers one message. Production code wraps an
// email service; tests use a fake.
type Sender interface {
	Send(ctx context.Context, to, subject, body string) error
}

// Message builds the reminder text for one loan.
func Message(o store.Overdue, now time.Time) string {
	days := fines.DaysLate(o.DueAt, now)
	return fmt.Sprintf("%q was due %d days ago. Fine so "+
		"far: %s.", o.Title, days,
		fines.FormatRupees(fines.Standard.Amount(days)))
}

// Run sends one reminder per overdue loan, one at a time.
// It stops at the first failure and reports how many went.
func Run(ctx context.Context, src Source, snd Sender,
	now time.Time) (int, error) {

	list, err := src.OverdueAt(ctx, now)
	if err != nil {
		return 0, fmt.Errorf("listing overdue loans: %w", err)
	}
	sent := 0
	for _, o := range list {
		err := snd.Send(ctx, o.Email, "Library book overdue",
			Message(o, now))
		if err != nil {
			return sent,
				fmt.Errorf("loan %d: %w", o.LoanID, err)
		}
		sent++
	}
	return sent, nil
}

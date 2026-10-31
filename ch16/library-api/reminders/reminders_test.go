package reminders

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/rohit-schoolabe/go-for-backend-jobs/ch16/library-api/store"
)

// fakeSource returns a fixed list instead of querying SQL.
type fakeSource struct {
	list []store.Overdue
	err  error
}

func (f fakeSource) OverdueAt(context.Context,
	time.Time) ([]store.Overdue, error) {
	return f.list, f.err
}

// fakeSender records messages and can fail on purpose.
type fakeSender struct {
	sent   []string
	failAt int // 1-based; 0 means never fail
}

func (f *fakeSender) Send(_ context.Context,
	to, _, body string) error {
	if len(f.sent)+1 == f.failAt {
		return errors.New("mail server down")
	}
	f.sent = append(f.sent, to+": "+body)
	return nil
}

var now = time.Date(2026, 2, 1, 6, 0, 0, 0, time.UTC)

func twoLate() fakeSource {
	due := now.AddDate(0, 0, -3)
	return fakeSource{list: []store.Overdue{
		{LoanID: 1, Email: "a@x.in", Title: "Optics",
			DueAt: due},
		{LoanID: 2, Email: "b@x.in", Title: "Go Basics",
			DueAt: due},
	}}
}

func TestRunSendsEveryReminder(t *testing.T) {
	snd := &fakeSender{}
	n, err := Run(context.Background(), twoLate(), snd, now)
	if err != nil || n != 2 {
		t.Fatalf("got %d, %v", n, err)
	}
	if !strings.Contains(snd.sent[0], "Rs 6.00") {
		t.Errorf("message: %s", snd.sent[0])
	}
}

func TestRunStopsOnSendError(t *testing.T) {
	snd := &fakeSender{failAt: 2}
	n, err := Run(context.Background(), twoLate(), snd, now)
	if err == nil || n != 1 {
		t.Fatalf("got %d, %v; want 1 and an error", n, err)
	}
}

func TestRunSourceError(t *testing.T) {
	src := fakeSource{err: errors.New("db locked")}
	_, err := Run(context.Background(), src,
		&fakeSender{}, now)
	if err == nil {
		t.Fatal("want the source error")
	}
}

package remind

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"golang.org/x/time/rate"

	"github.com/rohit-schoolabe/go-for-backend-jobs/ch16/library-api/store"
)

type manyLate int

func (n manyLate) OverdueAt(_ context.Context,
	now time.Time) ([]store.Overdue, error) {
	list := make([]store.Overdue, n)
	for i := range list {
		list[i] = store.Overdue{LoanID: int64(i),
			Email: fmt.Sprintf("m%d@x.in", i),
			Title: "Optics", DueAt: now.AddDate(0, 0, -2)}
	}
	return list, nil
}

// slowSender takes a while and records how many sends
// were in flight at once.
type slowSender struct {
	now, peak atomic.Int64
	failEvery int64
	calls     atomic.Int64
}

func (s *slowSender) Send(context.Context,
	string, string, string) error {
	n := s.now.Add(1)
	defer s.now.Add(-1)
	for {
		p := s.peak.Load()
		if n <= p || s.peak.CompareAndSwap(p, n) {
			break
		}
	}
	time.Sleep(5 * time.Millisecond)
	c := s.calls.Add(1)
	if s.failEvery > 0 && c%s.failEvery == 0 {
		return errors.New("bounced")
	}
	return nil
}

func TestPoolBoundsConcurrency(t *testing.T) {
	snd := &slowSender{failEvery: 10}
	p := Pool{Workers: 4}
	res, err := p.Run(context.Background(), manyLate(40),
		snd, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	if res.Sent != 36 || res.Failed != 4 {
		t.Errorf("got %+v, want 36 sent, 4 failed", res)
	}
	if peak := snd.peak.Load(); peak > 4 {
		t.Errorf("%d sends at once, limit is 4", peak)
	}
}

func TestPoolStopsOnCancel(t *testing.T) {
	ctx, cancel := context.WithTimeout(
		context.Background(), 50*time.Millisecond)
	defer cancel()
	p := Pool{Workers: 2,
		Limiter: rate.NewLimiter(rate.Limit(20), 1)}
	res, err := p.Run(ctx, manyLate(1000),
		&slowSender{}, time.Now())
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("err = %v", err)
	}
	if res.Sent > 10 {
		t.Errorf("sent %d after a 50ms deadline", res.Sent)
	}
}

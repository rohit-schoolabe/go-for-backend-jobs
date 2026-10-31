// Package cancel answers Chapter 17, exercise 2.
package cancel

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/rohit-schoolabe/go-for-backend-jobs/ch16/library-api/store"
	"github.com/rohit-schoolabe/go-for-backend-jobs/ch17/library-api/remind"
)

type forty struct{}

func (forty) OverdueAt(context.Context,
	time.Time) ([]store.Overdue, error) {
	list := make([]store.Overdue, 40)
	for i := range list {
		list[i] = store.Overdue{LoanID: int64(i)}
	}
	return list, nil
}

// cancelAfter cancels the run from inside the fifth send.
type cancelAfter struct {
	n      atomic.Int64
	cancel context.CancelFunc
}

func (c *cancelAfter) Send(context.Context,
	string, string, string) error {
	if c.n.Add(1) == 5 {
		c.cancel()
	}
	time.Sleep(time.Millisecond)
	return nil
}

func TestCancelMidRun(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	snd := &cancelAfter{cancel: cancel}
	p := remind.Pool{Workers: 2}
	res, err := p.Run(ctx, forty{}, snd, time.Now())
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("err = %v", err)
	}
	if res.Sent >= 40 {
		t.Errorf("sent all %d despite cancel", res.Sent)
	}
}

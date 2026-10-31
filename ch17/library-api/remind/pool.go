// Package remind sends overdue reminders concurrently,
// with a fixed number of workers and a rate limit.
package remind

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/time/rate"

	"github.com/rohit-schoolabe/go-for-backend-jobs/ch16/library-api/reminders"
	"github.com/rohit-schoolabe/go-for-backend-jobs/ch16/library-api/store"
)

// Pool is a worker pool. Limiter may be nil.
type Pool struct {
	Workers int
	Limiter *rate.Limiter
}

// Result counts what happened. One failed email does not
// stop the others.
type Result struct {
	Sent, Failed int64
}

// Run sends one reminder per overdue loan.
func (p Pool) Run(ctx context.Context, src reminders.Source,
	snd reminders.Sender, now time.Time) (Result, error) {

	list, err := src.OverdueAt(ctx, now)
	if err != nil {
		return Result{}, err
	}

	jobs := make(chan store.Overdue)
	var sent, failed atomic.Int64
	var wg sync.WaitGroup

	for range p.Workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for o := range jobs {
				if p.Limiter != nil {
					if p.Limiter.Wait(ctx) != nil {
						return // canceled while waiting
					}
				}
				err := snd.Send(ctx, o.Email,
					"Library book overdue",
					reminders.Message(o, now))
				if err != nil {
					failed.Add(1)
				} else {
					sent.Add(1)
				}
			}
		}()
	}

	// Feed the workers, but stop if the caller gives up.
feed:
	for _, o := range list {
		select {
		case jobs <- o:
		case <-ctx.Done():
			break feed
		}
	}
	close(jobs)
	wg.Wait()

	res := Result{Sent: sent.Load(), Failed: failed.Load()}
	return res, ctx.Err()
}

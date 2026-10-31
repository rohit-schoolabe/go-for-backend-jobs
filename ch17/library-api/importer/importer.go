// Package importer loads the old system's catalog export
// into Shelf as a three-stage pipeline.
package importer

import (
	"context"
	"io"
	"strings"
	"sync"
	"sync/atomic"

	"golang.org/x/sync/errgroup"

	"github.com/rohit-schoolabe/go-for-backend-jobs/ch13/catalogio"
	"github.com/rohit-schoolabe/go-for-backend-jobs/ch16/library-api/store"
)

// Inserter saves one book. *store.Store satisfies it.
type Inserter interface {
	CreateBook(ctx context.Context,
		b store.Book) (store.Book, error)
}

// Stats reports how the import went.
type Stats struct {
	Imported, Skipped int64
}

// Import decodes, cleans and inserts every book in r.
// Cleaning fans out to workers; inserting is a single
// goroutine, because SQLite allows one writer at a time.
func Import(ctx context.Context, r io.Reader, db Inserter,
	addedBy int64, workers int) (Stats, error) {

	g, ctx := errgroup.WithContext(ctx)
	raw := make(chan catalogio.Book)
	clean := make(chan store.Book)
	var stats Stats
	var skipped atomic.Int64

	// Stage 1: decode, one book at a time.
	g.Go(func() error {
		defer close(raw)
		return catalogio.ReadArray(r,
			func(b catalogio.Book) error {
				select {
				case raw <- b:
					return nil
				case <-ctx.Done():
					return ctx.Err()
				}
			})
	})

	// Stage 2: fan out to clean and validate.
	var cleaners sync.WaitGroup
	for range workers {
		cleaners.Add(1)
		g.Go(func() error {
			defer cleaners.Done()
			for b := range raw {
				sb, ok := tidy(b, addedBy)
				if !ok {
					skipped.Add(1)
					continue
				}
				select {
				case clean <- sb:
				case <-ctx.Done():
					return ctx.Err()
				}
			}
			return nil
		})
	}
	go func() {
		cleaners.Wait()
		close(clean)
	}()

	// Stage 3: fan in to one writer.
	g.Go(func() error {
		for b := range clean {
			if _, err := db.CreateBook(ctx, b); err != nil {
				return err
			}
			stats.Imported++
		}
		return nil
	})

	err := g.Wait()
	stats.Skipped = skipped.Load()
	return stats, err
}

// tidy trims the old system's padding and rejects rows
// that Shelf cannot store.
func tidy(b catalogio.Book, addedBy int64) (store.Book, bool) {
	t := strings.TrimSpace(b.Title)
	a := strings.TrimSpace(b.Author)
	if t == "" || a == "" || b.Year < 1450 || b.Year > 2100 {
		return store.Book{}, false
	}
	return store.Book{Title: t, Author: a, Year: b.Year,
		AddedBy: addedBy}, true
}

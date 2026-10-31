// Command library-api serves Shelf and runs its batch jobs:
//
//	library-api             serve the API
//	library-api remind      send overdue reminders
//	library-api import F    import the old catalog export
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"golang.org/x/time/rate"

	"github.com/rohit-schoolabe/go-for-backend-jobs/ch16/library-api/api"
	"github.com/rohit-schoolabe/go-for-backend-jobs/ch16/library-api/auth"
	"github.com/rohit-schoolabe/go-for-backend-jobs/ch16/library-api/store"
	"github.com/rohit-schoolabe/go-for-backend-jobs/ch17/library-api/importer"
	"github.com/rohit-schoolabe/go-for-backend-jobs/ch17/library-api/remind"
)

func main() {
	if err := run(); err != nil {
		slog.Error("failed", "err", err)
		os.Exit(1)
	}
}

// run holds the real main, so that its deferred calls
// finish before main calls os.Exit.
func run() error {
	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt, syscall.SIGTERM)
	defer stop()

	st, err := store.Open("library.db")
	if err != nil {
		return err
	}
	defer st.Close()

	switch cmd := arg(1); cmd {
	case "":
		return serve(ctx, st)
	case "remind":
		return sendReminders(ctx, st)
	case "import":
		return importFile(ctx, st, arg(2), arg(3))
	default:
		return errors.New("unknown command " + cmd)
	}
}

func arg(i int) string {
	if len(os.Args) > i {
		return os.Args[i]
	}
	return ""
}

func serve(ctx context.Context, st *store.Store) error {
	secret := os.Getenv("JWT_SECRET")
	if len(secret) < 32 {
		return errors.New("set JWT_SECRET to 32+ characters")
	}
	srv := &api.Server{Store: st,
		Tokens: auth.Tokens{Secret: []byte(secret),
			TTL: 2 * time.Hour},
		LoanDays: 14, Now: time.Now}
	addr := os.Getenv("ADDR") // e.g. ":9090"
	if addr == "" {
		addr = ":8080"
	}
	hs := &http.Server{
		Addr:              addr,
		Handler:           srv.Routes(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	errc := make(chan error, 1)
	go func() {
		slog.Info("listening", "addr", hs.Addr)
		errc <- hs.ListenAndServe()
	}()

	select {
	case err := <-errc:
		return err // could not start, e.g. port in use
	case <-ctx.Done():
	}

	slog.Info("shutting down, finishing open requests")
	shutCtx, cancel := context.WithTimeout(
		context.Background(), 10*time.Second)
	defer cancel()
	if err := hs.Shutdown(shutCtx); err != nil {
		return err
	}
	if err := <-errc; !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	slog.Info("stopped cleanly")
	return nil
}

// logSender stands in for an email service.
type logSender struct{}

func (logSender) Send(_ context.Context,
	to, subject, body string) error {
	slog.Info("email", "to", to, "subject", subject,
		"body", body)
	return nil
}

func sendReminders(ctx context.Context, st *store.Store) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Minute)
	defer cancel()
	p := remind.Pool{Workers: 8,
		Limiter: rate.NewLimiter(rate.Limit(10), 1)}
	res, err := p.Run(ctx, st, logSender{}, time.Now())
	slog.Info("reminders", "sent", res.Sent,
		"failed", res.Failed)
	return err
}

func importFile(ctx context.Context, st *store.Store,
	path, user string) error {
	uid, err := strconv.ParseInt(user, 10, 64)
	if err != nil {
		return errors.New("usage: import FILE USER_ID")
	}
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	stats, err := importer.Import(ctx, f, st, uid, 4)
	slog.Info("import", "imported", stats.Imported,
		"skipped", stats.Skipped)
	return err
}

// Command shelf is the library API, built to run in a
// container: settings from the environment, migrations
// on start, probes, and a clean stop on SIGTERM.
package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/rohit-schoolabe/go-for-backend-jobs/ch16/library-api/auth"
	"github.com/rohit-schoolabe/go-for-backend-jobs/ch19/shelf/api"
	"github.com/rohit-schoolabe/go-for-backend-jobs/ch19/shelf/config"
	"github.com/rohit-schoolabe/go-for-backend-jobs/ch19/shelf/store"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "healthcheck" {
		os.Exit(healthcheck())
	}
	if err := run(); err != nil {
		slog.Error("shelf", "err", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load(os.Getenv)
	if err != nil {
		return err
	}
	st, err := store.Open(cfg.DatabaseURL) // migrates
	if err != nil {
		return fmt.Errorf("database: %w", err)
	}
	defer st.Close()

	srv := &api.Server{
		Store: st,
		Tokens: auth.Tokens{
			Secret: []byte(cfg.JWTSecret), TTL: 2 * time.Hour,
		},
		LoanDays: cfg.LoanDays,
		Now:      time.Now,
	}
	hs := &http.Server{
		Addr:              cfg.Addr,
		Handler:           srv.Routes(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt, syscall.SIGTERM)
	defer stop()

	errc := make(chan error, 1)
	go func() { errc <- hs.ListenAndServe() }()
	slog.Info("listening", "addr", cfg.Addr)

	select {
	case err := <-errc:
		return err // could not bind, for example
	case <-ctx.Done():
	}

	slog.Info("SIGTERM: draining", "delay", cfg.DrainDelay)
	srv.Drain()
	time.Sleep(cfg.DrainDelay)

	sctx, cancel := context.WithTimeout(context.Background(),
		cfg.ShutdownTimeout)
	defer cancel()
	if err := hs.Shutdown(sctx); err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}
	if err := <-errc; !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	slog.Info("stopped cleanly")
	return nil
}

// healthcheck lets the container check itself. The final
// image has no shell and no curl, so the binary does it.
func healthcheck() int {
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080"
	}
	if strings.HasPrefix(addr, ":") {
		addr = "localhost" + addr
	}
	c := http.Client{Timeout: 2 * time.Second}
	resp, err := c.Get("http://" + addr + "/healthz")
	if err != nil {
		return 1
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 1
	}
	return 0
}

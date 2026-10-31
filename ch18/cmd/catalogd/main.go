// Command catalogd runs the internal catalog service.
package main

import (
	"context"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rohit-schoolabe/go-for-backend-jobs/ch16/library-api/auth"
	"github.com/rohit-schoolabe/go-for-backend-jobs/ch16/library-api/store"
	"github.com/rohit-schoolabe/go-for-backend-jobs/ch18/catalogsvc"
)

func main() {
	if err := run(); err != nil {
		slog.Error("catalogd", "err", err)
		os.Exit(1)
	}
}

func run() error {
	st, err := store.Open(env("CATALOG_DB", "library.db"))
	if err != nil {
		return err
	}
	defer st.Close()

	tokens := auth.Tokens{
		Secret: []byte(env("JWT_SECRET", "dev-only-secret")),
		TTL:    24 * time.Hour,
	}
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	srv := catalogsvc.New(st, tokens, log)

	lis, err := net.Listen("tcp", env("GRPC_ADDR", ":9090"))
	if err != nil {
		return err
	}
	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt, syscall.SIGTERM)
	defer stop()
	go func() {
		<-ctx.Done()
		done := make(chan struct{})
		go func() { srv.GracefulStop(); close(done) }()
		select {
		case <-done:
		case <-time.After(10 * time.Second):
			srv.Stop() // cut off streams that never end
		}
	}()
	log.Info("listening", "addr", lis.Addr().String())
	return srv.Serve(lis)
}

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

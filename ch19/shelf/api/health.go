package api

import (
	"context"
	"io"
	"net/http"
	"time"
)

// healthz answers "is the process alive?". It checks
// nothing else, so a slow database never gets the
// container restarted.
func (s *Server) healthz(w http.ResponseWriter,
	r *http.Request) {
	io.WriteString(w, "ok\n")
}

// readyz answers "should traffic come here?". It fails
// while shutting down and while the database is away.
func (s *Server) readyz(w http.ResponseWriter,
	r *http.Request) {
	if s.draining.Load() {
		http.Error(w, "draining", http.StatusServiceUnavailable)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(),
		time.Second)
	defer cancel()
	if err := s.Store.Ping(ctx); err != nil {
		http.Error(w, "database unreachable",
			http.StatusServiceUnavailable)
		return
	}
	io.WriteString(w, "ready\n")
}

// Drain makes readyz fail so the load balancer stops
// sending new requests. Call it when SIGTERM arrives.
func (s *Server) Drain() {
	s.draining.Store(true)
}

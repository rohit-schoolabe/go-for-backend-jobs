// Package api holds the HTTP handlers for the library API.
package api

import (
	"encoding/json"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/rohit-schoolabe/go-for-backend-jobs/ch16/library-api/auth"
	"github.com/rohit-schoolabe/go-for-backend-jobs/ch19/shelf/store"
)

// Server holds what every handler needs.
type Server struct {
	Store    *store.Store
	Tokens   auth.Tokens
	LoanDays int
	Now      func() time.Time // swapped in tests

	draining atomic.Bool // set once shutdown starts
}

// Routes maps each method and path to a handler.
func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /signup", s.signup)
	mux.HandleFunc("POST /login", s.login)
	mux.HandleFunc("GET /books", s.listBooks)
	mux.HandleFunc("GET /books/{id}", s.getBook)

	mux.Handle("POST /books", s.auth(s.createBook))
	mux.Handle("PUT /books/{id}", s.auth(s.updateBook))
	mux.Handle("DELETE /books/{id}", s.auth(s.deleteBook))
	mux.Handle("POST /books/{id}/loan", s.auth(s.borrow))
	mux.Handle("DELETE /books/{id}/loan", s.auth(s.giveBack))
	mux.Handle("GET /me/loans", s.auth(s.myLoans))

	// Probes skip the request log; they arrive every
	// few seconds and would drown it.
	root := http.NewServeMux()
	root.HandleFunc("GET /healthz", s.healthz)
	root.HandleFunc("GET /readyz", s.readyz)
	root.Handle("/", logRequests(mux))
	return root
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func fail(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

// decode reads a JSON body of at most 1 MB and
// rejects fields the API does not know.
func decode(r *http.Request, v any) error {
	body := http.MaxBytesReader(nil, r.Body, 1<<20)
	dec := json.NewDecoder(body)
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}

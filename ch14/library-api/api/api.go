// Package api holds the HTTP handlers for the library API.
package api

import (
	"encoding/json"
	"net/http"

	"github.com/rohit-schoolabe/go-for-backend-jobs/ch14/library-api/store"
)

// Server holds what every handler needs.
type Server struct {
	Store *store.Store
}

// Routes maps each method and path to a handler.
func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /books", s.listBooks)
	mux.HandleFunc("GET /books/{id}", s.getBook)
	mux.HandleFunc("POST /books", s.createBook)
	mux.HandleFunc("PUT /books/{id}", s.updateBook)
	mux.HandleFunc("DELETE /books/{id}", s.deleteBook)
	return mux
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

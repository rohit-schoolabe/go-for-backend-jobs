package api

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"
)

func logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			log.Printf("%s %s %v",
				r.Method, r.URL.Path, time.Since(start))
		})
}

type ctxKey struct{}

// auth lets a request through only with a valid token,
// and stores the user ID in the request context.
func (s *Server) auth(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			tok, ok := strings.CutPrefix(header, "Bearer ")
			if !ok || tok == "" {
				fail(w, http.StatusUnauthorized,
					"missing bearer token")
				return
			}
			uid, err := s.Tokens.Verify(tok)
			if err != nil {
				fail(w, http.StatusUnauthorized,
					"invalid or expired token")
				return
			}
			ctx := context.WithValue(r.Context(), ctxKey{}, uid)
			next(w, r.WithContext(ctx))
		})
}

// userID reads the ID that auth stored.
func userID(r *http.Request) int64 {
	id, _ := r.Context().Value(ctxKey{}).(int64)
	return id
}

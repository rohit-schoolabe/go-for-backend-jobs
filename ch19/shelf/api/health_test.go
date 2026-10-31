package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rohit-schoolabe/go-for-backend-jobs/ch19/shelf/store"
)

func TestProbes(t *testing.T) {
	st, err := store.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	s := &Server{Store: st}
	h := s.Routes()

	probe := func(path string) int {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest("GET", path, nil))
		return rec.Code
	}
	if probe("/healthz") != 200 || probe("/readyz") != 200 {
		t.Fatal("probes should pass on a fresh server")
	}

	s.Drain()
	if got := probe("/readyz"); got != 503 {
		t.Errorf("readyz while draining = %d, want 503", got)
	}
	if got := probe("/healthz"); got != 200 {
		t.Errorf("healthz while draining = %d, want 200", got)
	}

	st.Close() // the database goes away
	s.draining.Store(false)
	got := probe("/readyz")
	if got != http.StatusServiceUnavailable {
		t.Errorf("readyz with no database = %d", got)
	}
}

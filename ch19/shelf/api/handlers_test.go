package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/rohit-schoolabe/go-for-backend-jobs/ch16/library-api/auth"
)

// TestProtectedRoutes calls the router directly with a
// ResponseRecorder: no network, no port.
func TestProtectedRoutes(t *testing.T) {
	ts := newTestServer(t) // for a real user and token
	good := login(t, ts, "meena@college.in")

	other := auth.Tokens{
		Secret: []byte("another-secret-32-bytes-long-000"),
		TTL:    time.Hour,
	}
	forged, _ := other.Issue(1, "meena@college.in")

	tests := []struct {
		name   string
		header string
		want   int
	}{
		{"no header", "", http.StatusUnauthorized},
		{"not bearer", "Basic abc", http.StatusUnauthorized},
		{"garbage", "Bearer abc.def.ghi",
			http.StatusUnauthorized},
		{"wrong secret", "Bearer " + forged,
			http.StatusUnauthorized},
		{"valid", "Bearer " + good, http.StatusOK},
	}
	h := ts.Config.Handler
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/me/loans", nil)
			if tt.header != "" {
				req.Header.Set("Authorization", tt.header)
			}
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)
			if rec.Code != tt.want {
				t.Errorf("status %d, want %d; body %s",
					rec.Code, tt.want,
					strings.TrimSpace(rec.Body.String()))
			}
		})
	}
}

package api

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/rohit-schoolabe/go-for-backend-jobs/ch15/library-api/auth"
	"github.com/rohit-schoolabe/go-for-backend-jobs/ch15/library-api/store"
)

type obj = map[string]any

func newTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	log.SetOutput(io.Discard)
	st, err := store.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { st.Close() })

	s := &Server{
		Store: st,
		Tokens: auth.Tokens{
			Secret: []byte("test-secret-32-bytes-long-000000"),
			TTL:    time.Hour,
		},
		LoanDays: 14,
		Now: func() time.Time {
			return time.Date(2026, 1, 10, 9, 0, 0, 0, time.UTC)
		},
	}
	ts := httptest.NewServer(s.Routes())
	t.Cleanup(ts.Close)
	return ts
}

// call sends one request and decodes a JSON object reply.
func call(t *testing.T, ts *httptest.Server,
	method, path, token string, body any) (int, obj) {

	t.Helper()
	var buf bytes.Buffer
	if body != nil {
		json.NewEncoder(&buf).Encode(body)
	}
	req, _ := http.NewRequest(method, ts.URL+path, &buf)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	out := obj{}
	json.NewDecoder(res.Body).Decode(&out)
	return res.StatusCode, out
}

func login(
	t *testing.T, ts *httptest.Server, email string,
) string {
	t.Helper()
	cred := obj{"email": email, "password": "secret123"}
	code, _ := call(t, ts, "POST", "/signup", "", cred)
	if code != 201 {
		t.Fatalf("signup: got %d", code)
	}
	code, body := call(t, ts, "POST", "/login", "", cred)
	if code != 200 {
		t.Fatalf("login: got %d", code)
	}
	return body["token"].(string)
}

func TestAuthRules(t *testing.T) {
	ts := newTestServer(t)
	login(t, ts, "asha@example.com")
	book := obj{"title": "x", "author": "y", "year": 2000}

	tests := []struct {
		name, path, token string
		body              obj
		want              int
	}{
		{"duplicate signup", "/signup", "",
			obj{"email": "asha@example.com",
				"password": "secret123"}, 409},
		{"wrong password", "/login", "",
			obj{"email": "asha@example.com",
				"password": "wrongpass"}, 401},
		{"no token", "/books", "", book, 401},
		{"bad token", "/books", "garbage", book, 401},
	}
	for _, tt := range tests {
		code, _ := call(t, ts, "POST", tt.path, tt.token,
			tt.body)
		if code != tt.want {
			t.Errorf("%s: want %d, got %d",
				tt.name, tt.want, code)
		}
	}
}

func TestBookAndLoanFlow(t *testing.T) {
	ts := newTestServer(t)
	asha := login(t, ts, "asha@example.com")
	ravi := login(t, ts, "ravi@example.com")

	code, book := call(t, ts, "POST", "/books", asha, obj{
		"title": "Wings of Fire", "year": 1999,
		"author": "A. P. J. Abdul Kalam"})
	if code != 201 {
		t.Fatalf("create book: got %d", code)
	}
	path := "/books/" + strconv.Itoa(int(book["id"].(float64)))

	code, loan := call(t, ts, "POST", path+"/loan", ravi, nil)
	if code != 201 {
		t.Fatalf("borrow: got %d", code)
	}
	if loan["due_at"] != "2026-01-24T09:00:00Z" {
		t.Errorf("due_at: got %v", loan["due_at"])
	}

	steps := []struct {
		name, method, path, token string
		want                      int
	}{
		{"non-owner edit", "PUT", path, ravi, 403},
		{"double borrow", "POST", path + "/loan", asha, 409},
		{"delete on loan", "DELETE", path, asha, 409},
		{"return other's", "DELETE", path + "/loan", asha, 404},
		{"return", "DELETE", path + "/loan", ravi, 204},
		{"owner delete", "DELETE", path, asha, 204},
		{"get deleted", "GET", path, "", 404},
	}
	for _, st := range steps {
		code, _ := call(t, ts, st.method, st.path, st.token,
			nil)
		if code != st.want {
			t.Errorf("%s: want %d, got %d",
				st.name, st.want, code)
		}
	}
}

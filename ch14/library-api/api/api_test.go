package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/rohit-schoolabe/go-for-backend-jobs/ch14/library-api/store"
)

type obj = map[string]any

func newTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	st, err := store.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { st.Close() })
	ts := httptest.NewServer((&Server{Store: st}).Routes())
	t.Cleanup(ts.Close)
	return ts
}

func call(t *testing.T, ts *httptest.Server,
	method, path string, body any) (int, obj) {

	t.Helper()
	var buf bytes.Buffer
	if body != nil {
		json.NewEncoder(&buf).Encode(body)
	}
	req, _ := http.NewRequest(method, ts.URL+path, &buf)
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	out := obj{}
	json.NewDecoder(res.Body).Decode(&out)
	return res.StatusCode, out
}

func TestBookCRUD(t *testing.T) {
	ts := newTestServer(t)

	code, b := call(t, ts, "POST", "/books", obj{
		"title": "Wings of Fire", "year": 1999,
		"author": "A. P. J. Abdul Kalam"})
	if code != http.StatusCreated {
		t.Fatalf("create: got %d", code)
	}
	path := "/books/" + strconv.Itoa(int(b["id"].(float64)))

	steps := []struct {
		name, method, path string
		body               any
		want               int
	}{
		{"missing title", "POST", "/books",
			obj{"author": "x", "year": 2000}, 400},
		{"unknown field", "POST", "/books",
			obj{"title": "x", "author": "y", "year": 2000,
				"price": 10}, 400},
		{"get", "GET", path, nil, 200},
		{"bad id", "GET", "/books/abc", nil, 400},
		{"missing id", "GET", "/books/999", nil, 404},
		{"update", "PUT", path, obj{"title": "Wings of Fire",
			"author": "Kalam", "year": 1999}, 200},
		{"delete", "DELETE", path, nil, 204},
		{"get deleted", "GET", path, nil, 404},
	}
	for _, st := range steps {
		code, _ := call(t, ts, st.method, st.path, st.body)
		if code != st.want {
			t.Errorf("%s: want %d, got %d",
				st.name, st.want, code)
		}
	}
}

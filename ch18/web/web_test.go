package web

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	catalogv1 "github.com/rohit-schoolabe/go-for-backend-jobs/ch18/gen/shelf/catalog/v1"
)

// fakeCatalog answers GetBook without a network.
type fakeCatalog struct {
	catalogv1.CatalogServiceClient
}

func (fakeCatalog) GetBook(_ context.Context,
	req *catalogv1.GetBookRequest, _ ...grpc.CallOption,
) (*catalogv1.GetBookResponse, error) {
	if req.GetId() != 1 {
		return nil, status.Error(codes.NotFound, "no book")
	}
	return &catalogv1.GetBookResponse{Book: &catalogv1.Book{
		Id: 1, Title: "Godaan", Author: "Premchand",
		Availability: catalogv1.
			Availability_AVAILABILITY_ON_SHELF,
	}}, nil
}

func TestBookHandler(t *testing.T) {
	mux := http.NewServeMux()
	mux.Handle("GET /books/{id}", BookHandler(fakeCatalog{}))

	cases := []struct {
		path string
		code int
		body string
	}{
		{"/books/1", 200, `"availability":`},
		{"/books/7", 404, "no book"},
		{"/books/x", 400, "bad id"},
	}
	for _, tc := range cases {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", tc.path, nil)
		mux.ServeHTTP(rec, req)
		if rec.Code != tc.code {
			t.Errorf("%s: %d, want %d",
				tc.path, rec.Code, tc.code)
		}
		if !strings.Contains(rec.Body.String(), tc.body) {
			t.Errorf("%s: body %q", tc.path, rec.Body)
		}
	}
}

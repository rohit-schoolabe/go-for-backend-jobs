// Package web is the REST side: it answers HTTP requests
// by calling the catalog service over gRPC.
package web

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"

	catalogv1 "github.com/rohit-schoolabe/go-for-backend-jobs/ch18/gen/shelf/catalog/v1"
)

// BookHandler serves GET /books/{id}.
func BookHandler(
	c catalogv1.CatalogServiceClient,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			http.Error(w, "bad id", http.StatusBadRequest)
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(),
			2*time.Second)
		defer cancel()
		resp, err := c.GetBook(ctx,
			&catalogv1.GetBookRequest{Id: id})
		if err != nil {
			code := httpStatus(err)
			http.Error(w, status.Convert(err).Message(), code)
			return
		}
		out, err := protojson.Marshal(resp.GetBook())
		if err != nil {
			http.Error(w, "encoding", 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
	}
}

// httpStatus maps a gRPC status code to an HTTP one.
func httpStatus(err error) int {
	switch status.Code(err) {
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.NotFound:
		return http.StatusNotFound
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.DeadlineExceeded:
		return http.StatusGatewayTimeout
	case codes.Unavailable:
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}

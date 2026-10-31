// Command web is a small REST front end that reads books
// from catalogd over gRPC.
package main

import (
	"log"
	"net/http"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	catalogv1 "github.com/rohit-schoolabe/go-for-backend-jobs/ch18/gen/shelf/catalog/v1"
	"github.com/rohit-schoolabe/go-for-backend-jobs/ch18/web"
)

func main() {
	target := os.Getenv("CATALOG_ADDR")
	if target == "" {
		target = "localhost:9090"
	}
	// One connection, shared by every request. It is safe
	// for concurrent use and reconnects on its own.
	conn, err := grpc.NewClient(target,
		grpc.WithTransportCredentials(
			insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	client := catalogv1.NewCatalogServiceClient(conn)

	mux := http.NewServeMux()
	mux.Handle("GET /books/{id}", web.BookHandler(client))
	log.Fatal(http.ListenAndServe(":8081", mux))
}

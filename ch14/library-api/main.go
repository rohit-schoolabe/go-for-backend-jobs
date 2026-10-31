// Command library-api serves the Shelf library over HTTP.
package main

import (
	"log"
	"net/http"
	"time"

	"github.com/rohit-schoolabe/go-for-backend-jobs/ch14/library-api/api"
	"github.com/rohit-schoolabe/go-for-backend-jobs/ch14/library-api/store"
)

func main() {
	st, err := store.Open("library.db")
	if err != nil {
		log.Fatal(err)
	}
	defer st.Close()

	srv := &api.Server{Store: st}
	httpServer := &http.Server{
		Addr:              ":8080",
		Handler:           srv.Routes(),
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Println("library API on :8080")
	log.Fatal(httpServer.ListenAndServe())
}

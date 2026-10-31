// Command library-api serves the Shelf library over HTTP.
package main

import (
	"log"
	"net/http"
	_ "net/http/pprof" // adds /debug/pprof to DefaultServeMux
	"os"
	"time"

	"github.com/rohit-schoolabe/go-for-backend-jobs/ch16/library-api/api"
	"github.com/rohit-schoolabe/go-for-backend-jobs/ch16/library-api/auth"
	"github.com/rohit-schoolabe/go-for-backend-jobs/ch16/library-api/store"
)

func main() {
	secret := os.Getenv("JWT_SECRET")
	if len(secret) < 32 {
		log.Fatal("set JWT_SECRET to 32+ random characters")
	}
	st, err := store.Open("library.db")
	if err != nil {
		log.Fatal(err)
	}
	defer st.Close()

	srv := &api.Server{
		Store: st,
		Tokens: auth.Tokens{
			Secret: []byte(secret),
			TTL:    2 * time.Hour,
		},
		LoanDays: 14,
		Now:      time.Now,
	}
	httpServer := &http.Server{
		Addr:              ":8080",
		Handler:           srv.Routes(),
		ReadHeaderTimeout: 5 * time.Second,
	}
	// Profiles on a port only this machine can reach.
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	log.Println("library API on :8080")
	log.Fatal(httpServer.ListenAndServe())
}

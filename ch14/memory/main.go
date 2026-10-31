package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

var books = []book{{ID: 1, Title: "Go for Backend",
	Author: "Asha Nair"}}

type H map[string]any

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func listBooks(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, books)
}

func createBook(w http.ResponseWriter, r *http.Request) {
	var b book
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil || b.Title == "" || b.Author == "" {
		writeJSON(w, http.StatusBadRequest,
			H{"error": "title and author are required"})
		return
	}
	b.ID = len(books) + 1
	books = append(books, b) // not safe under load
	writeJSON(w, http.StatusCreated, b)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /books", listBooks)
	mux.HandleFunc("POST /books", createBook)
	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

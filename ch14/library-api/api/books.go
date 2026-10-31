package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/rohit-schoolabe/go-for-backend-jobs/ch14/library-api/store"
)

type bookInput struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Year   int    `json:"year"`
}

func (b bookInput) valid() bool {
	return strings.TrimSpace(b.Title) != "" &&
		strings.TrimSpace(b.Author) != "" &&
		b.Year > 0 && b.Year <= 2100
}

const badBook = "title, author and a sensible year are required"

func (s *Server) listBooks(
	w http.ResponseWriter, r *http.Request,
) {
	author := r.URL.Query().Get("author")
	books, err := s.Store.Books(r.Context(), author)
	if err != nil {
		fail(w, http.StatusInternalServerError,
			"could not list books")
		return
	}
	writeJSON(w, http.StatusOK, books)
}

// loadBook parses {id} and fetches the book. On failure it
// writes the error response itself and returns false.
func (s *Server) loadBook(
	w http.ResponseWriter, r *http.Request,
) (store.Book, bool) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		fail(w, http.StatusBadRequest, "id must be a number")
		return store.Book{}, false
	}
	b, err := s.Store.Book(r.Context(), id)
	if errors.Is(err, store.ErrNotFound) {
		fail(w, http.StatusNotFound, "book not found")
		return b, false
	}
	if err != nil {
		fail(w, http.StatusInternalServerError,
			"could not load book")
		return b, false
	}
	return b, true
}

func (s *Server) getBook(
	w http.ResponseWriter, r *http.Request,
) {
	if b, ok := s.loadBook(w, r); ok {
		writeJSON(w, http.StatusOK, b)
	}
}

func (s *Server) createBook(
	w http.ResponseWriter, r *http.Request,
) {
	var in bookInput
	if err := decode(r, &in); err != nil || !in.valid() {
		fail(w, http.StatusBadRequest, badBook)
		return
	}
	b := store.Book{Title: in.Title, Author: in.Author,
		Year: in.Year}
	b, err := s.Store.CreateBook(r.Context(), b)
	if err != nil {
		fail(w, http.StatusInternalServerError,
			"could not save book")
		return
	}
	writeJSON(w, http.StatusCreated, b)
}

func (s *Server) updateBook(
	w http.ResponseWriter, r *http.Request,
) {
	b, ok := s.loadBook(w, r)
	if !ok {
		return
	}
	var in bookInput
	if err := decode(r, &in); err != nil || !in.valid() {
		fail(w, http.StatusBadRequest, badBook)
		return
	}
	b.Title, b.Author, b.Year = in.Title, in.Author, in.Year
	if err := s.Store.UpdateBook(r.Context(), b); err != nil {
		fail(w, http.StatusInternalServerError,
			"could not update book")
		return
	}
	writeJSON(w, http.StatusOK, b)
}

func (s *Server) deleteBook(
	w http.ResponseWriter, r *http.Request,
) {
	b, ok := s.loadBook(w, r)
	if !ok {
		return
	}
	err := s.Store.DeleteBook(r.Context(), b.ID)
	if err != nil {
		fail(w, http.StatusInternalServerError,
			"could not delete book")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

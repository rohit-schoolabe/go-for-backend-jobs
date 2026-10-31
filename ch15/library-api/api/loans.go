package api

import (
	"errors"
	"net/http"

	"github.com/rohit-schoolabe/go-for-backend-jobs/ch15/library-api/store"
)

func (s *Server) borrow(
	w http.ResponseWriter, r *http.Request,
) {
	b, ok := s.loadBook(w, r)
	if !ok {
		return
	}
	now := s.Now().UTC()
	loan, err := s.Store.Borrow(r.Context(), b.ID, userID(r),
		now, s.LoanDays)
	if errors.Is(err, store.ErrOnLoan) {
		fail(w, http.StatusConflict, err.Error())
		return
	}
	if err != nil {
		fail(w, http.StatusInternalServerError,
			"could not borrow book")
		return
	}
	loan.BookTitle = b.Title
	writeJSON(w, http.StatusCreated, loan)
}

func (s *Server) giveBack(
	w http.ResponseWriter, r *http.Request,
) {
	b, ok := s.loadBook(w, r)
	if !ok {
		return
	}
	now := s.Now().UTC()
	err := s.Store.Return(r.Context(), b.ID, userID(r), now)
	if errors.Is(err, store.ErrNotFound) {
		fail(w, http.StatusNotFound,
			"you have no open loan on this book")
		return
	}
	if err != nil {
		fail(w, http.StatusInternalServerError,
			"could not return book")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) myLoans(
	w http.ResponseWriter, r *http.Request,
) {
	loans, err := s.Store.LoansFor(r.Context(), userID(r))
	if err != nil {
		fail(w, http.StatusInternalServerError,
			"could not list loans")
		return
	}
	writeJSON(w, http.StatusOK, loans)
}

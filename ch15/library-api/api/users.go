package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/rohit-schoolabe/go-for-backend-jobs/ch15/library-api/auth"
	"github.com/rohit-schoolabe/go-for-backend-jobs/ch15/library-api/store"
)

type credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *Server) signup(
	w http.ResponseWriter, r *http.Request,
) {
	var c credentials
	err := decode(r, &c)
	if err != nil || !strings.Contains(c.Email, "@") ||
		len(c.Password) < 8 {
		fail(w, http.StatusBadRequest,
			"need an email and a password of 8+ characters")
		return
	}
	hash, err := auth.HashPassword(c.Password)
	if err != nil {
		fail(w, http.StatusInternalServerError,
			"could not hash password")
		return
	}
	email := strings.ToLower(c.Email)
	u, err := s.Store.CreateUser(r.Context(), email, hash)
	if errors.Is(err, store.ErrDuplicate) {
		fail(w, http.StatusConflict, "email already registered")
		return
	}
	if err != nil {
		fail(w, http.StatusInternalServerError,
			"could not create user")
		return
	}
	writeJSON(w, http.StatusCreated, u)
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	var c credentials
	if err := decode(r, &c); err != nil {
		fail(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	email := strings.ToLower(c.Email)
	u, err := s.Store.UserByEmail(r.Context(), email)
	if err != nil && !errors.Is(err, store.ErrNotFound) {
		fail(w, http.StatusInternalServerError,
			"could not log in")
		return
	}
	if err != nil ||
		!auth.CheckPassword(c.Password, u.PasswordHash) {
		fail(w, http.StatusUnauthorized,
			"invalid email or password")
		return
	}
	tok, err := s.Tokens.Issue(u.ID, u.Email)
	if err != nil {
		fail(w, http.StatusInternalServerError,
			"could not issue token")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"token": tok})
}

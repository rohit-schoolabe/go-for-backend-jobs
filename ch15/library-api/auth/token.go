package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Tokens issues and checks signed login tokens.
type Tokens struct {
	Secret []byte
	TTL    time.Duration
}

type claims struct {
	UserID int64  `json:"uid"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// Issue returns a signed HS256 token for the user.
func (t Tokens) Issue(id int64, email string) (string, error) {
	now := time.Now()
	c := claims{
		UserID: id,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(t.TTL)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString(t.Secret)
}

// Verify checks signature, algorithm and expiry,
// and returns the user ID inside the token.
func (t Tokens) Verify(s string) (int64, error) {
	var c claims
	key := func(*jwt.Token) (any, error) {
		return t.Secret, nil
	}
	_, err := jwt.ParseWithClaims(s, &c, key,
		jwt.WithValidMethods([]string{"HS256"}))
	if err != nil {
		return 0, errors.New("invalid token")
	}
	return c.UserID, nil
}

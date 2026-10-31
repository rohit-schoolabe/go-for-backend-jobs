// Package auth hashes passwords and issues login tokens.
package auth

import "golang.org/x/crypto/bcrypt"

// HashPassword returns a salted bcrypt hash.
func HashPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword(
		[]byte(password), bcrypt.DefaultCost)
	return string(b), err
}

// CheckPassword reports whether password matches hash.
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hash), []byte(password))
	return err == nil
}

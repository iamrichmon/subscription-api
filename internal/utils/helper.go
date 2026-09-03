package utils

import (
	"errors"
	"net/mail"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// name normalization: trim spaces and reduce multiple spaces to single space
func NormalizeName(name string) string {
	name = strings.TrimSpace(name)
	return strings.Join(strings.Fields(name), " ")
}

// ErrEmailTaken is returned when an email is already registered.
var ErrEmailTaken = errors.New("Email already registered.")

// ErrInvalidCredentials is returned when the provided credentials are invalid.
var ErrInvalidCredentials = errors.New("invalid credentials")

// ErrInternalServerError is returned when an internal server error occurs.
var ErrInternalServerError = errors.New("internal server error")

// HashPlainPass hashes a plain password using bcrypt.
func HashPlainPass(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(hash), nil
}

// email normalization: trim spaces and convert to lowercase
func NormalizeEmail(email string) (string, error) {

	addr, err := mail.ParseAddress(strings.TrimSpace(email))

	if err != nil {
		return "", ErrInvalidCredentials
	}

	return strings.ToLower(addr.Address), nil
}

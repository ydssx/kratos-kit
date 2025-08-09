package util

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a plaintext password using bcrypt with a sensible default cost.
func HashPassword(plaintext string) (string, error) {
	if plaintext == "" {
		return "", errors.New("password cannot be empty")
	}
	// bcrypt.DefaultCost is a balanced choice for most environments
	hashed, err := bcrypt.GenerateFromPassword([]byte(plaintext), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

// ComparePassword compares a bcrypt hashed password with its possible plaintext equivalent.
// Returns nil on match, or an error otherwise.
func ComparePassword(hashed, plaintext string) error {
	if hashed == "" {
		return errors.New("empty hashed password")
	}
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plaintext))
}

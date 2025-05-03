package utils

import (
	"os"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	// Combine secret with password (optional step for extra salting)
	secretSaltedPassword := os.Getenv("SECRET")

	// Generate bcrypt hash
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(secretSaltedPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(passwordHash), nil
}

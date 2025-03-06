package models

import (
	"crypto/rand"
	"math/big"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_-+=<>?"

// GenerateRandomPassword generates a random password of a specified length
func GenerateRandomPassword(length int) string {
	password := make([]byte, length)
	for i := range password {
		// Generate a random index from the charset
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return ""
		}

		password[i] = charset[index.Int64()]
	}
	return string(password)
}

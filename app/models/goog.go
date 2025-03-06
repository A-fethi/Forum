package models

import (
	"database/sql"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// User model structure
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`  // Can be empty for Google login users
	GoogleID string `json:"google_id,omitempty"` // Can be set for Google login users
}

// CreateUser inserts a new user into the database
func CreateUser(db *sql.DB, user User) (int, error) {
	var userID int

	password1 := GenerateRandomPassword(12)
	password, _ := bcrypt.GenerateFromPassword([]byte(password1), 10)
	fmt.Println(password)
	// If the user has a Google ID, insert into the database with Google ID
	if user.GoogleID != "" {
		// Create a user with Google ID (no password)
		_, err := db.Exec(`
			INSERT INTO users (username, email,password, google_id)
			VALUES (?, ?,?, ?)`,
			user.Username, user.Email, password, user.GoogleID)
		if err != nil {
			return 0, fmt.Errorf("error creating Google user: %v", err)
		}
	} else {
		// Create a user with password (for normal users)
		_, err := db.Exec(`
			INSERT INTO users (username, email, password)
			VALUES (?, ?, ?)`,
			user.Username, user.Email, password)
		if err != nil {
			return 0, fmt.Errorf("error creating user: %v", err)
		}
	}

	// Retrieve the ID of the newly created user (whether Google login or normal registration)
	err := db.QueryRow(`
		SELECT id FROM users WHERE email = ? OR username = ?`,
		user.Email, user.Username).Scan(&userID)

	if err != nil {
		return 0, fmt.Errorf("error retrieving user ID: %v", err)
	}

	return userID, nil
}

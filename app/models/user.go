package models

import (
	"database/sql"
	"forum/app/config"
	"net/http"
	"regexp"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

type UserCredentials struct {
	Username string
	Email    string
	Password string
}

func VerifyPassword(hashedPassword, plainPassword string) bool {

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	if err != nil {

		return false
	}

	return true
}

func ValidUserName(name string) bool {
	if name == "" {
		return false
	}
	for _, char := range name {

		if !(unicode.IsLetter(char) || unicode.IsDigit(char)) {
			return false
		}
	}
	return true
}

func ValidEmail(email string) bool {
	valid := regexp.MustCompile(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`)
	return valid.MatchString(email)
}

func ValidatePassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	return true
}

func (User *UserCredentials) ValidInfo(resp http.ResponseWriter, db *sql.DB) (bool, string, int) {

	if !ValidUserName(User.Username) {
		config.Logger.Println("Failed Register Attempt: Username contains invalid characters.")
		return false, "username contains invalid characters", http.StatusBadRequest
	}

	var username string
	err := db.QueryRow("SELECT username FROM users WHERE username=?", User.Username).Scan(&username)
	if err != nil && err != sql.ErrNoRows {
		config.Logger.Println("Database error while checking username:", err)
		return false, "internal server error", http.StatusInternalServerError
	} else if err == nil {
		return false, "username already exists", http.StatusConflict
	}

	if !ValidEmail(User.Email) {
		return false, "invalid email format", http.StatusBadRequest
	}

	var email string
	err = db.QueryRow("SELECT email FROM users WHERE email=?", User.Email).Scan(&email)
	if err != nil && err != sql.ErrNoRows {
		config.Logger.Println("Database error while checking email:", err)
		return false, "internal server error", http.StatusInternalServerError
	} else if err == nil {
		return false, "email already registered", http.StatusConflict
	}

	if !ValidatePassword(User.Password) {
		return false, "password does not meet security criteria", http.StatusBadRequest
	}

	return true, "all inputs are valid", http.StatusOK
}

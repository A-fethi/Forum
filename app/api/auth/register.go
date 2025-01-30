package auth

import (
	"database/sql"
	"encoding/json"
	"forum/app/config"
	"forum/app/models"
	"forum/app/utils"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = 12

func Register(resp http.ResponseWriter, req *http.Request, db *sql.DB) {
	var potentialUser models.UserCredentials

	err := json.NewDecoder(req.Body).Decode(&potentialUser)
	if err != nil {
		config.Logger.Println("Register: Invalid request body:", err)
		models.SendErrorResponse(resp, http.StatusBadRequest, "Error: Invalid Request Format")
		return
	}

	valid, message, status := potentialUser.ValidInfo(resp, db)
	if !valid {
		models.SendErrorResponse(resp, status, message)
		return
	}

	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)", potentialUser.Email).Scan(&exists)
	if err != nil {
		config.Logger.Println("Register: Database error checking existing user:", err)
		models.SendErrorResponse(resp, http.StatusInternalServerError, "Error: Internal Server Error")
		return
	}
	if exists {
		config.Logger.Println("Register: Email already registered:", potentialUser.Email)
		models.SendErrorResponse(resp, http.StatusConflict, "Error: Email already registered")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(potentialUser.Password), bcryptCost)
	if err != nil {
		config.Logger.Println("Register: Error hashing password:", err)
		models.SendErrorResponse(resp, http.StatusInternalServerError, "Error: Internal Server Error")
		return
	}

	result, err := db.Exec(`
		INSERT INTO users (username, email, password) 
		VALUES (?, ?, ?)`,
		potentialUser.Username, potentialUser.Email, hashedPassword)
	if err != nil {
		config.Logger.Println("Register: Error inserting user into DB:", err)
		models.SendErrorResponse(resp, http.StatusInternalServerError, "Error: Internal Server Error")
		return
	}

	userID, err := result.LastInsertId()
	if err != nil {
		config.Logger.Println("Register: Error retrieving last inserted ID:", err)
		models.SendErrorResponse(resp, http.StatusInternalServerError, "Error: Internal Server Error")
		return
	}

	sessionToken, err := utils.ManageSession(db, int(userID), potentialUser.Username)
	if err != nil {
		config.Logger.Println("Register: Error managing session:", err)
		models.SendErrorResponse(resp, http.StatusInternalServerError, "Error: Internal Server Error")
		return
	}

	http.SetCookie(resp, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(2 * time.Hour),
	})

	config.Logger.Println("Register: Successful registration for user:", potentialUser.Username, "User ID:", userID)

	resp.WriteHeader(http.StatusCreated)
	// json.NewEncoder(resp).Encode(struct {
	// 	Message string `json:"message"`
	// }{
	// 	Message: "Registration Succes, logged In Automatically",
	// })
	// config.Templates.Exec(resp, "home.html", config.Login{
	// 	IsAuthenticated: true,
	// 	Username:        potentialUser.Username,
	// })

}

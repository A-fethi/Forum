package auth

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"forum/app/models"
	"forum/app/utils"
	"io/ioutil"
	"net/http"
)

func GoogleLogin(resp http.ResponseWriter, req *http.Request, db *sql.DB) {
	var requestData struct {
		IDToken string `json:"idToken"`
	}
	// Decode the ID token sent from the frontend
	if err := json.NewDecoder(req.Body).Decode(&requestData); err != nil {
		http.Error(resp, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Verify the ID token with Google API
	email, name, googleID, err := verifyGoogleIDToken(requestData.IDToken)
	if err != nil {
		http.Error(resp, "Google ID token verification failed", http.StatusUnauthorized)
		return
	}

	// Check if the user already exists in the database
	var user models.User
	err = db.QueryRow(`
		SELECT id, username, email, google_id FROM users WHERE google_id = ?`,
		googleID).Scan(&user.ID, &user.Username, &user.Email, &user.GoogleID)

	// If the user doesn't exist, create a new user
	if err == sql.ErrNoRows {
		user = models.User{
			Username: name,
			Email:    email, // Use email from the token
			GoogleID: googleID,
		}

		// Create the user
		_, err := models.CreateUser(db, user)
		if err != nil {
			http.Error(resp, "Error creating user", http.StatusInternalServerError)
			return
		}
	}

	// Generate session token and set the cookie
	sessionToken, err := utils.ManageSession(db, user.ID, user.Username)
	if err != nil {
		http.Error(resp, "Error generating session token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(resp, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Path:     "/",
		HttpOnly: true,
	})

	// Send response
	resp.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(resp).Encode(map[string]interface{}{
		"success": true,
		"user":    user,
		"message": "Google login successful",
	})
}

// Function to verify the Google ID token by calling Google's endpoint
func verifyGoogleIDToken(idToken string) (string, string, string, error) {
	url := fmt.Sprintf("https://www.googleapis.com/oauth2/v3/tokeninfo?id_token=%s", idToken)
	resp, err := http.Get(url)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to verify ID token: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return "", "", "", fmt.Errorf("Google verification failed: %s", string(body))
	}

	var payload struct {
		Email    string `json:"email"`
		GoogleID string `json:"sub"`
		Name     string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return "", "", "", fmt.Errorf("failed to parse Google response: %v", err)
	}

	return payload.Email, payload.Name, payload.GoogleID, nil
}

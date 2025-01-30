package posts

import (
	"database/sql"
	"encoding/json"
	"forum/app/api/auth"
	"forum/app/config"
	"forum/app/models"
	"forum/app/utils"
	"log"
	"net/http"
	"strings"
	"time"
)

func AddPost(resp http.ResponseWriter, req *http.Request, db *sql.DB) {
	// Ensure the user is authenticated before proceeding
	if !auth.SessionCheck(resp, req, db) {
		http.Error(resp, "User not authenticated", http.StatusUnauthorized)
		return
	}

	var request struct {
		Title       string   `json:"title"`
		PostContent string   `json:"content"`
		Categories  []string `json:"categories"`
	}

	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(resp, "400 - Invalid request body", http.StatusBadRequest)
		return
	}

	if err := utils.ValidatePost(request.Title, request.PostContent); err != nil {
		models.SendErrorResponse(resp, http.StatusBadRequest, "Error: Invalid Title/Post")
		return
	}

	sessionToken, err := utils.GetSessionToken(req)
	if err != nil || sessionToken == "" || !auth.SessionCheck(resp, req, db) {
		models.SendErrorResponse(resp, http.StatusUnauthorized, "Access: Unauthorized")
		// http.Error(resp, "User not authenticated", http.StatusUnauthorized)
		return
	}
	catCHECK := utils.CategoriesCheck(request.Categories)
	if !catCHECK {
		models.SendErrorResponse(resp, http.StatusBadRequest, "Error: Invalid Categories")
		return
	}
	_, username, err := utils.GetUsernameByToken(sessionToken, db)
	if err != nil {
		config.Logger.Println("Failed to get username:", err)
		models.SendErrorResponse(resp, http.StatusInternalServerError, "Error: Internal Server Error. Try later")
		// http.Error(resp, "Failed to get username", http.StatusInternalServerError)
		return
	}

	title := request.Title
	postContent := request.PostContent

	_, err = db.Exec(`
		INSERT INTO posts (username, title, content, categories, created_at)
		VALUES (?, ?, ?, ?, ?)`,
		username, title, postContent, strings.Join(request.Categories, " "), time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		config.Logger.Println("Failed to insert post:", err)
		models.SendErrorResponse(resp, http.StatusInternalServerError, "Error: Internal Server Error. Try later")
		// http.Error(resp, "Failed to insert post", http.StatusInternalServerError)
		return
	}

	var postID int
	err = db.QueryRow("SELECT id FROM posts WHERE username = ? ORDER BY created_at DESC LIMIT 1", username).Scan(&postID)
	if err == sql.ErrNoRows {
		config.Logger.Println("No posts found for the user.")
		// will never occur since we aready add the post above lol
		models.SendErrorResponse(resp, http.StatusNotFound, "No posts found for the user")
		return
	} else if err != nil {
		config.Logger.Println("Error retrieving last post ID:", err)
		models.SendErrorResponse(resp, http.StatusInternalServerError, "Error: Internal Server Error. Try later")
		// models.SendErrorResponse(resp, http.StatusInternalServerError, "Error retrieving user data")
		return
	}

	post := models.Post{
		Username:   username,
		ID:         postID,
		Title:      title,
		Content:    postContent,
		Categories: strings.Join(request.Categories, " "),
		CreatedAt:  time.Now().Format("2006-01-02 15:04:05"),
		Likes:      0,
		Dislikes:   0,
	}

	config.Logger.Printf("Post created successfully, postID: %d", postID)
	resp.WriteHeader(http.StatusCreated)
	resp.Header().Set("Content-Type", "application/json")
	json.NewEncoder(resp).Encode(post)
}

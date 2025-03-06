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
	var (
		postID int
	)
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
		return
	}

	title := request.Title
	postContent := request.PostContent

	err = db.QueryRow(`
		INSERT INTO posts (username, title, content, categories, created_at)
		VALUES (?, ?, ?, ?, ?)
		RETURNING id`,
		username, title, postContent, strings.Join(request.Categories, " "), time.Now()).Scan(&postID)
	if err != nil {
		config.Logger.Println("Failed to insert post:", err)
		models.SendErrorResponse(resp, http.StatusInternalServerError, "Error: Internal Server Error. Try later")
		return
	}

	post := models.Post{
		Username:   username,
		ID:         postID,
		Title:      title,
		Content:    postContent,
		Categories: strings.Join(request.Categories, " "),
		CreatedAt:  utils.TimeAgo(time.Now()),
		Likes:      0,
		Dislikes:   0,
	}

	config.Logger.Printf("Post created successfully, postID: %d", postID)
	resp.WriteHeader(http.StatusCreated)
	resp.Header().Set("Content-Type", "application/json")
	json.NewEncoder(resp).Encode(post)
}

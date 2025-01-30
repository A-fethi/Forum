package comments

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"forum/app/api/auth"
	"forum/app/config"
	"forum/app/models"
	"forum/app/utils"
)

func AddComment(resp http.ResponseWriter, req *http.Request, db *sql.DB) {
	var Comment struct {
		Content string `json:"content"`
		Post_id int    `json:"post_id"`
	}

	err := json.NewDecoder(req.Body).Decode(&Comment)
	if err != nil {
		config.Logger.Printf("Error decoding request body: %v\n", err)
		models.SendErrorResponse(resp, http.StatusBadRequest, "Error: Invalid Content")
		// http.Error(resp, "400 - invalid request body", http.StatusBadRequest)
		return
	}

	session_token, err := utils.GetSessionToken(req)
	if err != nil || session_token == "" || !auth.SessionCheck(resp, req, db) {
		config.Logger.Println("User not authenticated: ", err)
		models.SendErrorResponse(resp, http.StatusUnauthorized, "Access: Unauthorized")
		// http.Error(resp, "401 - Unauthorized", http.StatusUnauthorized)
		return
	}

	if Comment.Content == "" || Comment.Post_id == 0 {
		config.Logger.Println("Comment content, username cannot be empty")
		models.SendErrorResponse(resp, http.StatusBadRequest, "Error: Invalid Content")

		// http.Error(resp, "400 - invalid request body", http.StatusBadRequest)
		return
	}

	var username string
	var userID int
	userID, username, err = utils.GetUsernameByToken(session_token, db)
	if err != nil {
		models.SendErrorResponse(resp, http.StatusInternalServerError, "Error: Internal Server Error. Try later")
		// http.Error(resp, "Failed to get username", http.StatusInternalServerError)
		return
	}
	config.Logger.Println("Username: ", username, "UserID: ", userID, "Comment: ", Comment)
	_, err = db.Exec(`
		INSERT INTO comments (post_id, user_id, author, content, created_at)
		VALUES (?, ?, ?, ?, ?)`,
		Comment.Post_id, userID, username, Comment.Content, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		models.SendErrorResponse(resp, http.StatusInternalServerError, "Error: Internal Server Error. Try later")
		// http.Error(resp, "Failed to insert comment", http.StatusInternalServerError)
		return
	}

	var CommentID int
	err = db.QueryRow("SELECT id FROM comments WHERE author = ? ORDER BY created_at DESC LIMIT 1", username).Scan(&CommentID)
	if err == sql.ErrNoRows {
		config.Logger.Println("No posts found for the user.") // will never occure too
		models.SendErrorResponse(resp, http.StatusNotFound, "No posts found for the user")
		return
	} else if err != nil {
		config.Logger.Println("Error retrieving last post ID:", err)
		models.SendErrorResponse(resp, http.StatusInternalServerError, "Error: Internal Server Error. Try later")
		// models.SendErrorResponse(resp, http.StatusInternalServerError, "Error retrieving user data")
		return
	}
	config.Logger.Println("Comment add: ", Comment)
	comm := models.Comment{
		ID:        CommentID,
		Content:   Comment.Content,
		Username:  username,
		PostID:    Comment.Post_id,
		CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
		Likes:     0,
		Dislikes:  0,
	}
	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(http.StatusCreated)
	json.NewEncoder(resp).Encode(comm)
}

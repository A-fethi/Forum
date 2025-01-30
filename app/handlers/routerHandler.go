package handlers

import (
	"database/sql"
	"forum/app/api/auth"
	"forum/app/api/comments"
	"forum/app/api/posts"
	"forum/app/api/reactions"
	"forum/app/models"
	"net/http"
	"strings"
)

func Router(resp http.ResponseWriter, req *http.Request, db *sql.DB) {

	path := strings.Split(req.URL.Path[5:], "/")
	switch path[0] {
	case "auth":
		auth.Authentication(resp, req, db)
	case "posts":
		if req.Method == http.MethodGet {
			data := posts.GetPosts(resp, req, db)
			if len(data) == 0 {
				resp.Header().Set("Content-Type", "application/json")
				resp.Write([]byte("[]"))
			} else {
				resp.Header().Set("Content-Type", "application/json")
				resp.Write(data)
			}

		} else if req.Method == http.MethodPost {
			posts.AddPost(resp, req, db)
		} else {
			models.SendErrorResponse(resp, http.StatusMethodNotAllowed, "Error: Method not allowed")
			return
		}
	case "comments":
		if req.Method == http.MethodPost {
			comments.AddComment(resp, req, db)
		} else {
			models.SendErrorResponse(resp, http.StatusMethodNotAllowed, "Error: Method not allowed")
			return
		}
	case "reactions":
		reactions.AddReaction(resp, req, db)
	default:
		models.SendErrorResponse(resp, http.StatusNotFound, "Page Not Found")
		return
	}
}

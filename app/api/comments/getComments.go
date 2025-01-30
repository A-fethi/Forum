package comments

import (
	"database/sql"
	"forum/app/models"
)

// func GetComments(resp http.ResponseWriter, req *http.Request, db *sql.DB) {

// }

func GetComments(postID int, db *sql.DB) ([]models.Comment, error) {
	// Get post ID from URL
	// postID := req.URL.Path[len("/api/comments/"):]
	rows, err := db.Query("SELECT id, author, content, created_at, likes, dislikes FROM comments WHERE post_id = ? ORDER BY created_at DESC", postID)
	if err != nil {

		return nil, err
	}
	defer rows.Close()

	var commentss []models.Comment
	for rows.Next() {
		var comments models.Comment
		if err := rows.Scan(&comments.ID, &comments.Username, &comments.Content, &comments.CreatedAt, &comments.Likes, &comments.Dislikes); err != nil {

			return nil, err
		}

		commentss = append(commentss, comments)
	}
	if err := rows.Err(); err != nil {

		return nil, err
	}

	// Prepare data for template
	return commentss, nil
}

package forum

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
)

func ToggleInteraction(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./database/database.db")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	userIDStr := r.FormValue("user_id")
	commentIDStr := r.FormValue("comment_id")
	action := r.FormValue("action")

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil || (action != "like" && action != "dislike") {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var existingID int
	var existingAction string
	query := `SELECT id, action FROM user_interactions WHERE user_id = ? AND comment_id = ?`
	err = db.QueryRow(query, userID, commentID).Scan(&existingID, &existingAction)

	if err == sql.ErrNoRows {
		_, err = db.Exec(
			`INSERT INTO user_interactions (user_id, comment_id, action) VALUES (?, ?, ?)`,
			userID, commentID, action,
		)
		if err != nil {
			log.Printf("Error adding interaction: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if action == "like" {
			db.Exec(`UPDATE comments SET likes = likes + 1 WHERE id = ?`, commentID)
		} else {
			db.Exec(`UPDATE comments SET dislikes = dislikes + 1 WHERE id = ?`, commentID)
		}
	} else if err == nil {
		if existingAction == action {
			if action == "like" {
				db.Exec(`UPDATE comments SET likes = CASE WHEN likes > 0 THEN likes - 1 ELSE likes END WHERE id = ?`, commentID)
			} else {
				db.Exec(`UPDATE comments SET dislikes = CASE WHEN dislikes > 0 THEN dislikes - 1 ELSE likes END WHERE id = ?`, commentID)
			}
			db.Exec(`DELETE FROM user_interactions WHERE id = ?`, existingID)
		} else {
			db.Exec(`DELETE FROM user_interactions WHERE id = ?`, existingID)
			db.Exec(
				`INSERT INTO user_interactions (user_id, comment_id, action) VALUES (?, ?, ?)`,
				userID, commentID, action,
			)
			if action == "like" {
				db.Exec(`UPDATE comments SET likes = likes + 1, dislikes = CASE WHEN dislikes > 0 THEN dislikes - 1 ELSE dislikes END WHERE id = ?`, commentID)
			} else {
				db.Exec(`UPDATE comments SET dislikes = dislikes + 1, likes = CASE WHEN likes > 0 THEN likes - 1 ELSE likes END WHERE id = ?`, commentID)
			}
		}
	} else {
		log.Printf("Error querying interactions: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/comments", http.StatusSeeOther)
}

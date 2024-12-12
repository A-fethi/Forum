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
	commentID, err2 := strconv.Atoi(commentIDStr)

	if err != nil || err2 != nil || (action != "like" && action != "dislike") {
		log.Fatal(err)
		log.Fatal(err2)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var existingID int
	query := `SELECT id FROM user_interactions WHERE user_id = ? AND comment_id = ? AND action = ?`
	err = db.QueryRow(query, userID, commentID, action).Scan(&existingID)

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
			_, err = db.Exec(`UPDATE comments SET likes = likes + 1 WHERE id = ?`, commentID)
		} else {
			_, err = db.Exec(`UPDATE comments SET dislikes = dislikes + 1 WHERE id = ?`, commentID)
		}
	} else if err == nil {
		_, err = db.Exec(
			`DELETE FROM user_interactions WHERE id = ?`,
			existingID,
		)
		if err != nil {
			log.Printf("Error removing interaction: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if action == "like" {
			_, err = db.Exec(`UPDATE comments SET likes = likes - 1 WHERE id = ?`, commentID)
		} else {
			_, err = db.Exec(`UPDATE comments SET dislikes = dislikes - 1 WHERE id = ?`, commentID)
		}
	} else {
		log.Printf("Error querying interactions: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/posts", http.StatusSeeOther)
}

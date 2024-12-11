package forum

import (
	"database/sql"
	"log"
	"net/http"
)

func HandleInteraction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	action := r.FormValue("action")
	// commentIDStr := r.FormValue("comment_id")
	// commentID, err := strconv.Atoi(commentIDStr)
	// if err != nil {
	// 	http.Error(w, "Invalid comment ID", http.StatusBadRequest)
	// 	return
	// }
	db, err := sql.Open("sqlite3", "./database/database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	switch action {
	case "Like":
		_, err = db.Exec("UPDATE comments SET likes = likes +1")
	case "Dislike":
		_, err = db.Exec("UPDATE comments SET dislikes = dislikes + 1")
	default:
		http.Error(w, "Invalid action", http.StatusBadRequest)
		return
	}

	if err != nil {
		log.Fatal(err)
		http.Error(w, "Failed to update count", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/posts", http.StatusSeeOther)
}

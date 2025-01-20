package forum

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

type Comment struct {
	ID       int
	Content  string
	Likes    int
	Dislikes int
}

type PageData struct {
	Comments []Comment
}

func HandleComments(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./database/database.db")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	if r.Method == http.MethodGet {
		displayComments(w, db)
		return
	}

	if r.Method == http.MethodPost {
		handleCommentRequest(w, r, db)
		// http.Redirect(w, r, "/comments", http.StatusSeeOther)
		return
	}

	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
}

// func handleCommentRequest(w http.ResponseWriter, r *http.Request, db *sql.DB) {
// 	comment := r.FormValue("comment")
	// reply := r.FormValue("reply")
	// commentIDStr := r.FormValue("comment_id")

	// if comment != "" {
	// 	newComment := Comment{Content: comment}
	// 	_, err := insertSqlComment(db, newComment)
	// 	if err != nil {
	// 		log.Printf("Error inserting comment: %v", err)
	// 		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	// 	}
	// }
	// if comment != "" {
	// 	newComment := Comment{Content: comment}
	// 	result, err := insertSqlComment(db, newComment)
	// 	if err != nil {
	// 		log.Printf("Error inserting comment: %v", err)
	// 		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	// 		return
	// 	}

	// 	id, _ := result.LastInsertId()
	// 	newComment.ID = int(id)

	// 	w.Header().Set("Content-Type", "application/json")
	// 	json.NewEncoder(w).Encode(newComment)
	// 	return
	// }

	// if reply != "" {
	// 	commentID, err := strconv.Atoi(commentIDStr)
	// 	if err != nil {
	// 		log.Printf("Invalid comment ID: %v", err)
	// 		http.Error(w, "Bad Request", http.StatusBadRequest)
	// 		return
	// 	}
	// 	newReply := Reply{Content: reply}
	// 	_, err = insertSqlReply(db, newReply, commentID)
	// 	if err != nil {
	// 		log.Printf("Error inserting reply: %v", err)
	// 		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	// 	}
	// }
// }

func handleCommentRequest(w http.ResponseWriter, r *http.Request, db *sql.DB) {
    comment := r.FormValue("comment")
    if comment != "" {
        // Use the existing insertSqlComment function
        newComment := Comment{Content: comment}
        result, err := insertSqlComment(db, newComment)
        if err != nil {
            log.Printf("Error inserting comment: %v", err)
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }

        // Get the ID of the newly inserted comment
        id, err := result.LastInsertId()
        if err != nil {
            log.Printf("Error getting last insert ID: %v", err)
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }

        // Fetch the full details of the newly inserted comment
        var insertedComment Comment
        err = db.QueryRow("SELECT id, content, likes, dislikes FROM comments WHERE id = ?", id).Scan(
            &insertedComment.ID, &insertedComment.Content, &insertedComment.Likes, &insertedComment.Dislikes,
        )
        if err != nil {
            log.Printf("Error fetching new comment: %v", err)
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }

        // Return the new comment as JSON
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(insertedComment)
        return
    }

    http.Error(w, "Bad Request", http.StatusBadRequest)
}

func displayComments(w http.ResponseWriter, db *sql.DB) {
	comments, err := getAllComments(db)
	if err != nil {
		log.Printf("Error fetching comments: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	pageData := PageData{Comments: comments}
	tmpl, err := template.ParseFiles("templates/comments.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, pageData)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func getAllComments(db *sql.DB) ([]Comment, error) {
	rows, err := db.Query("SELECT id, content, likes, dislikes FROM comments ORDER BY created_at ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		if err := rows.Scan(&comment.ID, &comment.Content, &comment.Likes, &comment.Dislikes); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

package forum

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
)

type Comment struct {
	ID      int
	Author  int
	Content string
}

type Reply struct {
	ID      int
	Author  int
	Content string
}

type PageData struct {
	Comments []Comment
	Replies  []Reply
}

func HandlePosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		HandleError(w, http.StatusMethodNotAllowed)
		return
	}

	var comments []Comment
	var replies []Reply
	comment := r.FormValue("comment")
	reply := r.FormValue("reply")

	if comment != "" {
		newComment := Comment{
			ID:      5,
			Content: comment,
		}

		db, err := sql.Open("sqlite3", "./database/database.db")
		if err != nil {
			log.Fatal(err)
			// HandleError(w, http.StatusInternalServerError)
			return
		}
		defer db.Close()

		_, err = insertSqlComment(db, newComment)
		if err != nil {
			log.Fatal(err)
			// HandleError(w, http.StatusInternalServerError)
			return
		}

		comments, err = getAllComments(db)
		if err != nil {
			log.Fatal(err)
			// HandleError(w, http.StatusInternalServerError)
			return
		}
	} 

	if reply != "" {
		newReply := Reply{
			ID:      5,
			Content: reply,
		}

		db, err := sql.Open("sqlite3", "./database/database.db")
		if err != nil {
			log.Fatal(err)
			// HandleError(w, http.StatusInternalServerError)
			return
		}
		defer db.Close()

		_, err = insertSqlReply(db, newReply)
		if err != nil {
			log.Fatal(err)
			// HandleError(w, http.StatusInternalServerError)
			return
		}

		replies, err = getAllReplies(db)
		if err != nil {
			log.Fatal(err)
			// HandleError(w, http.StatusInternalServerError)
			return
		}
	}

	pageData := PageData{
		Comments: comments,
		Replies: replies,
	}

	tmpl, err := template.ParseFiles("templates/posts.html")
	if err != nil {
		log.Fatal(err)
		// HandleError(w, http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, pageData)
	if err != nil {
		log.Fatal(err)
		// HandleError(w, http.StatusInternalServerError)
		return
	}
}

func getAllComments(db *sql.DB) ([]Comment, error) {
	row, err := db.Query("SELECT id, content FROM comments ORDER BY created_at DESC")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	var result []Comment
	for row.Next() {
		var comment Comment
		if err := row.Scan(&comment.ID, &comment.Content); err != nil {
			log.Fatal(err)
		}
		result = append(result, comment)
	}
	return result, nil
}

func getAllReplies(db *sql.DB) ([]Reply, error) {
	row, err := db.Query("SELECT id, content FROM replies ORDER BY created_at DESC")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	var result []Reply
	for row.Next() {
		var reply Reply
		if err := row.Scan(&reply.ID, &reply.Content); err != nil {
			log.Fatal(err)
		}
		result = append(result, reply)
	}
	return result, nil
}
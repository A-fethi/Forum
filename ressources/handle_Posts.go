package forum

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type Comment struct {
	ID       int
	Author   int
	Content  string
	Likes    int
	Dislikes int
	Replies  []Reply
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
		commentIDStr := r.FormValue("comment_id")
		commentID, err := strconv.Atoi(commentIDStr)
		if err != nil {
			HandleError(w, http.StatusBadRequest)
			return
		}

		newReply := Reply{
			Content: reply,
		}

		db, err := sql.Open("sqlite3", "./database/database.db")
		if err != nil {
			log.Fatal(err)
			return
		}
		defer db.Close()

		_, err = insertSqlReply(db, newReply, commentID)
		if err != nil {
			log.Fatal(err)
			return
		}

		comments, err = getAllComments(db)
		if err != nil {
			log.Fatal(err)
			return
		}
	}

	pageData := PageData{
		Comments: comments,
		Replies:  replies,
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
	row, err := db.Query("SELECT id, content, likes, dislikes FROM comments ORDER BY created_at ASC")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	var comments []Comment
	for row.Next() {
		var comment Comment
		if err := row.Scan(&comment.ID, &comment.Content, &comment.Likes, &comment.Dislikes); err != nil {
			log.Fatal(err)
		}
		comment.Replies, err = getRepliesForComment(db, comment.ID)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func getRepliesForComment(db *sql.DB, commentID int) ([]Reply, error) {
	row, err := db.Query("SELECT id, content FROM replies WHERE comment_id = ? ORDER BY created_at ASC", commentID)
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	var replies []Reply
	for row.Next() {
		var reply Reply
		if err := row.Scan(&reply.ID, &reply.Content); err != nil {
			log.Fatal(err)
		}
		replies = append(replies, reply)
	}
	return replies, nil
}

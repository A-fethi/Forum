package forum

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
)

type Comment struct {
	ID      int
	Author  string
	Content string
}

type Reply struct {
	ID      int
	Author  string
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

	comment := r.FormValue("comment")
	reply := r.FormValue("reply")

	newComment := Comment{
		Author:  "comment",
		Content: comment,
	}

	newReply := Reply{
		Author:  "reply",
		Content: reply,
	}

	db, err := sql.Open("sqlite3", "./database/database.db")
	if err != nil {
		HandleError(w, http.StatusInternalServerError)
		return
	}
	defer db.Close()

	_, err = insertSql(db, newComment, newReply)
	if err != nil {
		log.Fatal(err)
		HandleError(w, http.StatusInternalServerError)
		return
	}

	comments, replies, err := getAllComments(db)
	if err != nil {
		HandleError(w, http.StatusInternalServerError)
		return
	}

	pageData := PageData{
		Comments: comments,
		Replies: replies,
	}

	tmpl, err := template.ParseFiles("templates/posts.html")
	if err != nil {
		HandleError(w, http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, pageData)
	if err != nil {
		HandleError(w, http.StatusInternalServerError)
		return
	}
}

func getAllComments(db *sql.DB) ([]Comment, []Reply, error) {
	row, err := db.Query("SELECT author, content FROM comments")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	var result []Comment
	for row.Next() {
		var comment Comment
		if err := row.Scan(&comment.Author, &comment.Content); err != nil {
			log.Fatal(err)
		}
	}

	rows, err := db.Query("SELECT author, content FROM comments")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var replies []Reply
	for row.Next() {
		var reply Reply
		if err := row.Scan(&reply.Author, &reply.Content); err != nil {
			log.Fatal(err)
		}
		replies = append(replies, reply)
	}

	return result, replies, nil
}

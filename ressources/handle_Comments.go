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
	Content  string
	Likes    int
	Dislikes int
	Replies  []Reply
}

type Reply struct {
	ID      int
	Content string
}

type PageData struct {
	Comments []Comment
}

func HandlePosts(w http.ResponseWriter, r *http.Request) {
	// Open the database connection
	db, err := sql.Open("sqlite3", "./database/database.db")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	if r.Method == http.MethodGet {
		// Handle GET requests to display posts
		displayPosts(w, db)
		return
	}

	if r.Method == http.MethodPost {
		// Handle POST requests for comments or replies
		handlePostRequest(w, r, db)
		http.Redirect(w, r, "/posts", http.StatusSeeOther)
		return
	}

	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
}

func handlePostRequest(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	comment := r.FormValue("comment")
	reply := r.FormValue("reply")
	commentIDStr := r.FormValue("comment_id")

	if comment != "" {
		// Insert a new comment
		newComment := Comment{Content: comment}
		_, err := insertSqlComment(db, newComment)
		if err != nil {
			log.Printf("Error inserting comment: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}

	if reply != "" {
		// Validate and insert a new reply
		commentID, err := strconv.Atoi(commentIDStr)
		if err != nil {
			log.Printf("Invalid comment ID: %v", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		newReply := Reply{Content: reply}
		_, err = insertSqlReply(db, newReply, commentID)
		if err != nil {
			log.Printf("Error inserting reply: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

func displayPosts(w http.ResponseWriter, db *sql.DB) {
	// Fetch all comments and their associated replies
	comments, err := getAllComments(db)
	if err != nil {
		log.Printf("Error fetching comments: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Render the template with the fetched data
	pageData := PageData{Comments: comments}
	tmpl, err := template.ParseFiles("templates/posts.html")
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
		comment.Replies, err = getRepliesForComment(db, comment.ID)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func getRepliesForComment(db *sql.DB, commentID int) ([]Reply, error) {
	rows, err := db.Query("SELECT id, content FROM replies WHERE comment_id = ? ORDER BY created_at ASC", commentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var replies []Reply
	for rows.Next() {
		var reply Reply
		if err := rows.Scan(&reply.ID, &reply.Content); err != nil {
			return nil, err
		}
		replies = append(replies, reply)
	}
	return replies, nil
}

package forum

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
)


func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "StatusMethodNotAllowed", http.StatusMethodNotAllowed)
		// HandleError(w, http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/" {
		http.Error(w, "w, StatusNotFound", http.StatusNotFound)
		// HandleError(w, http.StatusNotFound)
		return
	}


	db, err := sql.Open("sqlite3", "./database/database.db")
	if err != nil {
		log.Fatal(err)
		// HandleError(w, http.StatusInternalServerError)
		return
	}
	defer db.Close()

	comments, err := getAllComments(db)
	if err != nil {
		log.Fatal(err)
		// HandleError(w, http.StatusInternalServerError)
		return
	}

	pageData := PageData{
		Comments: comments,
	}

	tmpl, err := template.ParseFiles("templates/posts.html")
	if err != nil {
		log.Fatal(err)
		http.Error(w, "StatusInternalServerError", http.StatusInternalServerError)
		// HandleError(w, http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, pageData)
	if err != nil {
		log.Fatal(err)
		http.Error(w, "StatusInternalServerError", http.StatusInternalServerError)
		// HandleError(w, http.StatusInternalServerError)
		return
	}
}

package forum

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
)

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")
		if email == "" || password == "" {
			HandleError(w, http.StatusBadRequest)
			return
		}

		db, err := sql.Open("sqlite3", "./database/database.db")
		if err != nil {
			HandleError(w, http.StatusInternalServerError)
			log.Fatal(err)
			return
		}
		defer db.Close()

		if r.Method == http.MethodPost {
			insertSql := `INSERT INTO users (email, password) VALUES (?, ?)`
			_, err := db.Exec(insertSql, email, password)
			if err != nil {
				HandleError(w, http.StatusInternalServerError)
				log.Fatal(err)
				return
			}
			tmpl, err := template.ParseFiles("templates/index.html")
			if err != nil {
				HandleError(w, http.StatusInternalServerError)
				log.Fatal(err)
				return
			}

			err = tmpl.Execute(w, nil)
			if err != nil {
				HandleError(w, http.StatusInternalServerError)
				log.Fatal(err)
				return
			}
		}
	} else {
		tmpl, err := template.ParseFiles("templates/login.html")
		if err != nil {
			HandleError(w, http.StatusInternalServerError)
			log.Fatal(err)
			return
		}

		err = tmpl.Execute(w, nil)
		if err != nil {
			HandleError(w, http.StatusInternalServerError)
			log.Fatal(err)
			return
		}
	}
}

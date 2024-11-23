package forum

import (
	"html/template"
	"net/http"
)

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		HandleError(w, http.StatusMethodNotAllowed)
		return
	}

	tmpl, err := template.ParseFiles("templates/login.html")
	if err != nil {
		HandleError(w, http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		HandleError(w, http.StatusInternalServerError)
		return
	}
}

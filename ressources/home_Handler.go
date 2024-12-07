package forum

import (
	"fmt"
	"html/template"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		HandleError(w, http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/" {
		HandleError(w, http.StatusNotFound)
		return
	}

	comment := r.FormValue("comment")
	fmt.Println("comment", comment)

	tmpl, err := template.ParseFiles("templates/posts.html")
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

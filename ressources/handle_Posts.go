package forum

import (
	"fmt"
	"html/template"
	"net/http"
)

type Comment struct {
	ID      int
	Author  string
	Content string
}

func HandlePosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		HandleError(w, http.StatusMethodNotAllowed)
		return
	}

	test := r.FormValue("comment")
	fmt.Println("test", test)

	data := &Comment{
		Content: test,
	}

	tmpl, err := template.ParseFiles("templates/test.html")
	if err != nil {
		HandleError(w, http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		HandleError(w, http.StatusInternalServerError)
		return
	}
}

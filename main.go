package main

import (
	"fmt"
	forum "forum/ressources"
	"net/http"
)

func main() {
	http.HandleFunc("/static/", forum.HandleStatic)
	http.HandleFunc("/", forum.HomeHandler)
	http.HandleFunc("/login", forum.HandleLogin)
	http.HandleFunc("/signup", forum.HandleSignup)
	fmt.Println("server listening on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

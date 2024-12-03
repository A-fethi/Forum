package main

import (
	"database/sql"
	"fmt"
	forum "forum/ressources"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./database/database.db")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer db.Close()

	_, err = CreateTable(db)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("All tables was created successfully.")

	http.HandleFunc("/static/", forum.HandleStatic)
	http.HandleFunc("/", forum.HomeHandler)
	http.HandleFunc("/login", forum.HandleLogin)
	http.HandleFunc("/signup", forum.HandleSignup)
	fmt.Println("server listening on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

package main

import "database/sql"

func CreateTable(db *sql.DB) (sql.Result, error) {
	sql := `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT NOT NULL,
		password TEXT NOT NULL
	);
	CREATE TABLE IF NOT EXISTS likes (
		likeId INTEGER PRIMARY KEY AUTOINCREMENT,
		userId INTEGER,
		postId INTEGER,
		like INTEGER
	);
	CREATE TABLE IF NOT EXISTS comments (
    	id INTEGER PRIMARY KEY AUTOINCREMENT,
    	post_id INTEGER,
    	author TEXT,
    	content TEXT,
    	FOREIGN KEY(post_id) REFERENCES posts(id)
);`

	return db.Exec(sql)
}

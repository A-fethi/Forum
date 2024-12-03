package main

import "database/sql"

func CreateTable(db *sql.DB) (sql.Result, error) {
	sql := `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY,
		email TEXT NOT NULL,
		password TEXT NOT NULL
	);`

	return db.Exec(sql)
}

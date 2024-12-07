package main

import "database/sql"

func insertSql(db *sql.DB, c *Comment) (int64, error) {
	sql := `INSERT INTO comments (Author, Content) VALUES (?, ?);`
	result, err := db.Exec(sql, c.Author, c.Content)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

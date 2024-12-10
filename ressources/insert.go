package forum

import "database/sql"

func insertSqlComment(db *sql.DB, c Comment) (int64, error) {
	sql := `INSERT INTO comments (user_id, Content) VALUES (?, ?);`
	result, err := db.Exec(sql, c.Author, c.Content)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func insertSqlReply(db *sql.DB, r Reply) (int64, error) {
	sql := `INSERT INTO replies (user_id, Content) VALUES (?, ?);`
	result, err := db.Exec(sql, r.Author, r.Content)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

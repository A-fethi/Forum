package forum

import "database/sql"

func insertSqlComment(db *sql.DB, c Comment) (int64, error) {
	sql := `INSERT INTO comments (user_id, Content) VALUES (?, ?);`
	result, err := db.Exec(sql, c.ID, c.Content)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func insertSqlReply(db *sql.DB, reply Reply, commentID int) (sql.Result, error) {
	stmt, err := db.Prepare("INSERT INTO replies (content, comment_id) VALUES (?, ?)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return stmt.Exec(reply.Content, commentID)
}

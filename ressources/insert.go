package forum

import "database/sql"

func insertSqlComment(db *sql.DB, c Comment) (sql.Result, error) {
	stmt, err := db.Prepare("INSERT INTO comments (user_id, Content) VALUES (?, ?)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return stmt.Exec(c.ID, c.Content)
}

func insertSqlReply(db *sql.DB, reply Reply, commentID int) (sql.Result, error) {
	stmt, err := db.Prepare("INSERT INTO replies (content, comment_id) VALUES (?, ?)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return stmt.Exec(reply.Content, commentID)
}

package forum

import "database/sql"

func insertSqlComment(db *sql.DB, c Comment) (sql.Result, error) {
    stmt, err := db.Prepare("INSERT INTO comments (content, likes, dislikes) VALUES (?, 0, 0)")
    if err != nil {
        return nil, err
    }
    defer stmt.Close()

    return stmt.Exec(c.Content)
}


// func insertSqlReply(db *sql.DB, reply Reply, commentID int) (sql.Result, error) {
// 	stmt, err := db.Prepare("INSERT INTO replies (content, comment_id) VALUES (?, ?)")
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer stmt.Close()

// 	return stmt.Exec(reply.Content, commentID)
// }

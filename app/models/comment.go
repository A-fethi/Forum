package models

type Comment struct {
	ID        int
	Content   string
	Username  string
	PostID    int
	CreatedAt string
	Likes     int // Number of likes for this comment
	Dislikes  int // Number of dislikes for this comment

}

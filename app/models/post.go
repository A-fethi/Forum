package models

type Post struct {
	ID         int
	Username   string
	Title      string
	Content    string
	Categories string
	Comments   []Comment
	CreatedAt  string
	Likes      int
	Dislikes   int
}

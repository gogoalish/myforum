package models

type Post struct {
	ID      int
	UserID  int
	Creator string
	Title   string
	Content string
}

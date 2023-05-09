package models

type Comment struct {
	ID            int
	PostID        int
	UserID        int
	Content       string
	ParentID      int
	Creator       string
	Replies       []*Comment
	LikesCount    int
	DislikesCount int
}

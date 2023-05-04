package models

type Post struct {
	ID            int
	UserID        int
	Title         string
	Content       string
	Creator       string
	CmntCount     int
	LikesCount    int
	DislikesCount int
	Comments      []*Comment
}

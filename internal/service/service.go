package service

import (
	"forum/internal/repository"
)

type Service struct {
	Users
	Posts
	Comments
}

func New(r *repository.Repo) *Service {
	return &Service{
		Users:    &UserService{r.Users},
		Posts:    &PostService{r.Posts},
		Comments: &CommentService{r.Comments},
	}
}

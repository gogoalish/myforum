package service

import (
	"forum/internal/models"
	"forum/internal/repository"
)

type Comments interface {
	Create(c *models.Comment) error
	Fetch(postID int) ([]*models.Comment, error)
}

type CommentService struct {
	repo repository.Comments
}

func (s *CommentService) Create(c *models.Comment) error {
	return s.repo.InsertComment(c)
}

func (s *CommentService) Fetch(postID int) ([]*models.Comment, error) {
	return s.repo.CommentsByPostId(postID)
}

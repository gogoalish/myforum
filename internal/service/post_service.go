package service

import (
	"forum/internal/models"
	"forum/internal/repository"
)

type PostService struct {
	repo repository.Posts
}

type Posts interface {
	GetAll() ([]*models.Post, error)
	Create(UserID int, title, content string) (int, error)
	GetById(id int) (*models.Post, error)
}

func (s *PostService) GetAll() ([]*models.Post, error) {
	posts, err := s.repo.FetchPosts()
	return posts, err
}

func (s *PostService) Create(UserID int, title, content string) (int, error) {
	return s.repo.InsertPost(UserID, title, content)
}

func (s *PostService) GetById(id int) (*models.Post, error) {
	return s.repo.PostById(id)
}

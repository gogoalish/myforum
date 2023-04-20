package service

import (
	"forum/internal/models"
	"forum/internal/repository"
)

type PostService struct {
	repo repository.Posts
}

type Posts interface {
	All() ([]*models.Post, error)
	Create(UserID int, title, content string) error
}

func (p *PostService) All() ([]*models.Post, error) {
	posts, err := p.repo.All()
	return posts, err
}

func (p *PostService) Create(UserID int, title, content string) error {
	err := p.repo.Create(UserID, title, content)
	return err
}

package service

import (
	"forum/internal/models"
	"forum/internal/repository"
)

type PostService struct {
	repo     repository.Posts
	cmntRepo repository.Comments
}

type Posts interface {
	GetAll() ([]*models.Post, error)
	Create(UserID int, title, content string) (int, error)
	GetById(id int) (*models.Post, error)
}

func (s *PostService) GetAll() ([]*models.Post, error) {
	posts, err := s.repo.FetchPosts()
	for _, post := range posts {
		post.CmntCount, err = s.cmntRepo.CountCommentsByPostId(post.ID)
		if err != nil {
			return nil, err
		}
	}
	return posts, err
}

func (s *PostService) Create(UserID int, title, content string) (int, error) {
	return s.repo.InsertPost(UserID, title, content)
}

func (s *PostService) GetById(id int) (*models.Post, error) {
	post, err := s.repo.PostById(id)
	if err != nil {
		return nil, err
	}
	post.CmntCount, err = s.cmntRepo.CountCommentsByPostId(id)
	if err != nil {
		return nil, err
	}
	return post, nil
}

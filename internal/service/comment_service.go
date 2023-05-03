package service

import (
	"errors"

	"forum/internal/models"
	"forum/internal/repository"
)

type Comments interface {
	Create(c *models.Comment) error
	Fetch(postID int) ([]*models.Comment, error)
	Count(postID int) (int, error)
}

type CommentService struct {
	repo repository.Comments
}

func (s *CommentService) Create(c *models.Comment) error {
	return s.repo.InsertComment(c)
}

func (s *CommentService) Fetch(postID int) ([]*models.Comment, error) {
	comments, err := s.repo.CommentsByPostId(postID)
	if err != nil {
		return comments, err
	}
	err = s.Recur(comments)
	if err != nil {
		return nil, err
	}
	return comments, err
}

func (s *CommentService) Recur(comments []*models.Comment) (err error) {
	for _, comment := range comments {
		comment.Replies, err = s.repo.RepliesByParent(comment.ID)
		if err != nil && !errors.Is(err, models.ErrNoRecord) {
			return err
		}
		if !errors.Is(err, models.ErrNoRecord) {
			err = s.Recur(comment.Replies)
			if err != nil && !errors.Is(err, models.ErrNoRecord) {
				return err
			}
		}
	}
	return nil
}

func (s *CommentService) Count(postID int) (int, error) {
	return s.repo.CountCommentsByPostId(postID)
}

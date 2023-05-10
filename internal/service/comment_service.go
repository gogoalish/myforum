package service

import (
	"errors"
	"fmt"

	"forum/internal/models"
	"forum/internal/repository"
)

type Comments interface {
	Create(c *models.Comment) error
	Fetch(postID int) ([]*models.Comment, error)
	Count(postID int) (int, error)
	GetByID(comID int) (*models.Comment, error)
	React(comID, userID int, received string) error
	CountLikes(postID int) (int, error)
	CountDislikes(postID int) (int, error)
}

var ErrInvalidParent = errors.New("invalid parent")

type CommentService struct {
	repo repository.Comments
}

func (s *CommentService) Create(c *models.Comment) error {
	count, err := s.repo.CountAllComments()
	if err != nil {
		return err
	}
	if c.ParentID > count {
		fmt.Println(c.ParentID, count)
		return ErrInvalidParent
	}
	return s.repo.InsertComment(c)
}

func (s *CommentService) GetByID(comID int) (*models.Comment, error) {
	return s.repo.CommentById(comID)
}

func (s *CommentService) Fetch(postID int) ([]*models.Comment, error) {
	comments, err := s.repo.CommentsByPostId(postID)
	if err != nil {
		return comments, err
	}
	err = s.GetReplies(comments)
	if err != nil {
		return nil, err
	}
	return comments, err
}

func (s *CommentService) GetReplies(comments []*models.Comment) (err error) {
	for _, comment := range comments {
		comment.LikesCount, err = s.CountLikes(comment.ID)
		if err != nil {
			return err
		}
		comment.DislikesCount, err = s.CountDislikes(comment.ID)
		if err != nil {
			return err
		}
		comment.Replies, err = s.repo.RepliesByParent(comment.ID)
		if err != nil && !errors.Is(err, models.ErrNoRecord) {
			return err
		}
		if !errors.Is(err, models.ErrNoRecord) {
			err = s.GetReplies(comment.Replies)
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

func (s *CommentService) React(comID, userID int, received string) error {
	reaction, err := s.repo.ReactionByUserId(comID, userID)
	if err != nil && !errors.Is(err, models.ErrNoRecord) {
		return err
	}
	switch reaction {
	case "":
		err = s.repo.InsertReaction(comID, userID, received)
		if err != nil {
			return err
		}
	case received:
		err = s.repo.RemoveReaction(comID, userID)
		if err != nil {
			return err
		}
	default:
		err = s.repo.UpdateReaction(comID, userID, received)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *CommentService) CountLikes(comID int) (int, error) {
	likes, err := s.repo.LikesByCommentId(comID)
	return len(likes), err
}

func (s *CommentService) CountDislikes(comID int) (int, error) {
	dislikes, err := s.repo.DislikesByCommentId(comID)
	return len(dislikes), err
}

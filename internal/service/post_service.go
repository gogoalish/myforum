package service

import (
	"errors"

	"forum/internal/models"
	"forum/internal/repository"
)

type PostService struct {
	repo     repository.Posts
	cmntRepo repository.Comments
}

type Posts interface {
	GetAll() ([]*models.Post, error)
	Create(p *models.Post) (int, error)
	GetById(id int) (*models.Post, error)
	React(postID, userID int, reaction string) error
	CountLikes(postID int) (int, error)
	CountDislikes(postID int) (int, error)
}

func (s *PostService) GetAll() ([]*models.Post, error) {
	posts, err := s.repo.FetchPosts()
	for _, post := range posts {
		post.CmntCount, err = s.cmntRepo.CountCommentsByPostId(post.ID)
		if err != nil {
			return nil, err
		}
		post.LikesCount, err = s.CountLikes(post.ID)
		if err != nil {
			return nil, err
		}
		post.DislikesCount, err = s.CountDislikes(post.ID)
		if err != nil {
			return nil, err
		}
	}
	return posts, err
}

func (s *PostService) Create(p *models.Post) (int, error) {
	return s.repo.InsertPost(p)
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
	post.LikesCount, err = s.CountLikes(post.ID)
	if err != nil {
		return nil, err
	}
	post.DislikesCount, err = s.CountDislikes(post.ID)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (s *PostService) React(postID, userID int, r string) error {
	reaction, err := s.repo.ReactionByUserId(postID, userID)
	if err != nil && !errors.Is(err, models.ErrNoRecord) {
		return err
	}
	switch reaction {
	case "":
		err = s.repo.InsertReaction(postID, userID, r)
		if err != nil {
			return err
		}
	case r:
		err = s.repo.RemoveReaction(postID, userID)
		if err != nil {
			return err
		}
	default:
		err = s.repo.UpdateReaction(postID, userID, r)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *PostService) CountLikes(postID int) (int, error) {
	likes, err := s.repo.LikesByPostId(postID)
	return len(likes), err
}

func (s *PostService) CountDislikes(postID int) (int, error) {
	dislikes, err := s.repo.DislikesByPostId(postID)
	return len(dislikes), err
}

package repository

import (
	"database/sql"
	"errors"

	"forum/internal/models"
)

type Comments interface {
	InsertComment(c *models.Comment) error
	CommentsByPostId(postID int) ([]*models.Comment, error)
	CountCommentsByPostId(postID int) (int, error)
	RepliesByParent(parentID int) ([]*models.Comment, error)
}

type CommentRepo struct {
	*sql.DB
}

func (r *CommentRepo) InsertComment(c *models.Comment) error {
	query := `INSERT INTO comments
	VALUES(NULL, ?, ?, ?, ?)`
	_, err := r.Exec(query, c.PostID, c.UserID, c.Content, c.ParentID)
	if err != nil {
		return err
	}
	return nil
}

func (r *CommentRepo) CommentsByPostId(postID int) ([]*models.Comment, error) {
	query := `SELECT *, (SELECT name FROM users WHERE comments.user_id = users.id) FROM comments
	WHERE post_id = ? AND parent_id = 0`
	comments := []*models.Comment{}
	rows, err := r.Query(query, postID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}
	for rows.Next() {
		c := &models.Comment{}
		err = rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Content, &c.ParentID, &c.Creator)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, nil
}

func (r *CommentRepo) RepliesByParent(parentID int) ([]*models.Comment, error) {
	query := `SELECT *, (SELECT name FROM users WHERE comments.user_id = users.id) FROM comments
	WHERE parent_id = ?`
	replies := []*models.Comment{}
	rows, err := r.Query(query, parentID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}
	for rows.Next() {
		c := &models.Comment{}
		err = rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Content, &c.ParentID, &c.Creator)
		if err != nil {
			return nil, err
		}
		replies = append(replies, c)
	}
	return replies, nil
}

func (r *CommentRepo) CountCommentsByPostId(postID int) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM comments WHERE post_id = ?`
	err := r.QueryRow(query, postID).Scan(&count)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrNoRecord
		}
		return 0, err
	}
	return count, nil
}

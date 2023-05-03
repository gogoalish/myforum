package repository

import (
	"database/sql"
	"errors"

	"forum/internal/models"
)

type Comments interface {
	InsertComment(c *models.Comment) error
	CommentsByPostId(postID int) ([]*models.Comment, error)
}

type CommentRepo struct {
	*sql.DB
}

func (r *CommentRepo) InsertComment(c *models.Comment) error {
	query := `INSERT INTO comments
	VALUES(NULL, ?, ?, ?)`
	_, err := r.Exec(query, c.PostID, c.UserID, c.Content)
	if err != nil {
		return err
	}
	return nil
}

func (r *CommentRepo) CommentsByPostId(postID int) ([]*models.Comment, error) {
	query := `SELECT *, (SELECT name FROM users WHERE comments.user_id = users.id) FROM comments
	WHERE post_id = ?`
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
		err = rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Content, &c.Creator)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	return comments, nil
}

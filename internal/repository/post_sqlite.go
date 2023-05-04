package repository

import (
	"database/sql"
	"errors"

	"forum/internal/models"
)

type PostRepo struct {
	*sql.DB
}

type Posts interface {
	InsertPost(p *models.Post) (int, error)
	FetchPosts() ([]*models.Post, error)
	PostById(id int) (*models.Post, error)
	InsertReaction(postID, userID int, reaction string) error
	ReactionByUserId(postID, userID int) (string, error)
	RemoveReaction(postID, userID int) error
	UpdateReaction(postID, userID int, reaction string) error
	LikesByPostId(postID int) ([]*models.Reaction, error)
	DislikesByPostId(postID int) ([]*models.Reaction, error)
}

func (r *PostRepo) InsertPost(p *models.Post) (int, error) {
	query := `INSERT INTO posts
	VALUES(NULL, ?, ?, ?)`
	res, err := r.Exec(query, p.UserID, p.Title, p.Content)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (r *PostRepo) PostById(id int) (*models.Post, error) {
	query := `SELECT *, (SELECT name FROM users WHERE users.id = posts.user_id) FROM posts WHERE id = ?`
	p := &models.Post{}
	err := r.QueryRow(query, id).Scan(&p.ID, &p.UserID, &p.Title, &p.Content, &p.Creator)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}
	return p, nil
}

func (r *PostRepo) FetchPosts() ([]*models.Post, error) {
	posts := []*models.Post{}
	query := `SELECT *, (
		SELECT name FROM users WHERE users.id = posts.user_id
		)
	FROM posts`
	rows, err := r.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		p := &models.Post{}
		err := rows.Scan(&p.ID, &p.UserID, &p.Title, &p.Content, &p.Creator)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, err
}

func (r *PostRepo) InsertReaction(postID, userID int, reaction string) error {
	query := `INSERT INTO reactions (post_id, user_id, type) VALUES(?, ?, ?)`
	_, err := r.Exec(query, postID, userID, reaction)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostRepo) ReactionByUserId(postID, userID int) (string, error) {
	query := `SELECT type FROM reactions WHERE post_id = ? AND user_id = ?`
	var reaction string
	err := r.QueryRow(query, postID, userID).Scan(&reaction)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return reaction, models.ErrNoRecord
		}
		return reaction, err
	}
	return reaction, nil
}

func (r *PostRepo) RemoveReaction(postID, userID int) error {
	query := "DELETE FROM reactions WHERE post_id = ? AND user_id = ?"
	_, err := r.Exec(query, postID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostRepo) UpdateReaction(postID, userID int, reaction string) error {
	query := `UPDATE reactions
	SET type = ?
	WHERE post_id = ? AND user_id = ?`
	_, err := r.Exec(query, reaction, postID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostRepo) LikesByPostId(postID int) ([]*models.Reaction, error) {
	query := `SELECT *
	FROM reactions WHERE post_id = ? AND type = "like"`
	rows, err := r.Query(query, postID)
	if err != nil {
		return nil, err
	}
	likes := []*models.Reaction{}
	defer rows.Close()
	for rows.Next() {
		r := &models.Reaction{}
		err := rows.Scan(&r.ID, &r.PostID, &r.CommentID, &r.UserID, &r.Type)
		if err != nil {
			return nil, err
		}
		likes = append(likes, r)
	}
	return likes, err
}

func (r *PostRepo) DislikesByPostId(postID int) ([]*models.Reaction, error) {
	query := `SELECT *
	FROM reactions WHERE post_id = ? AND type = "dislike"`
	rows, err := r.Query(query, postID)
	if err != nil {
		return nil, err
	}
	dislikes := []*models.Reaction{}
	defer rows.Close()
	for rows.Next() {
		r := &models.Reaction{}
		err := rows.Scan(&r.ID, &r.PostID, &r.CommentID, &r.UserID, &r.Type)
		if err != nil {
			return nil, err
		}
		dislikes = append(dislikes, r)
	}
	return dislikes, err
}

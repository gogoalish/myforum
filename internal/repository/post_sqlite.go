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
	InsertPost(UserID int, title, content string) (int, error)
	FetchPosts() ([]*models.Post, error)
	PostById(id int) (*models.Post, error)
}

func (r *PostRepo) InsertPost(userID int, title, content string) (int, error) {
	query := `INSERT INTO posts
	VALUES(NULL, ?, ?, ?)`
	res, err := r.Exec(query, userID, title, content)
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

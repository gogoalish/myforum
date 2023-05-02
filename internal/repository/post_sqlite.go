package repository

import (
	"database/sql"
	"errors"

	"forum/internal/models"
)

type PostRepo struct {
	DB *sql.DB
}

type Posts interface {
	Insert(UserID int, title, content string) (int, error)
	GetAll() ([]*models.Post, error)
	PostById(id int) (*models.Post, error)
}

func (r *PostRepo) Insert(userID int, title, content string) (int, error) {
	query := `INSERT INTO posts
	VALUES(NULL, ?, ?, ?)`
	res, err := r.DB.Exec(query, userID, title, content)
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
	err := r.DB.QueryRow(query, id).Scan(&p.ID, &p.UserID, &p.Title, &p.Content, &p.Creator)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}
	return p, nil
}

func (r *PostRepo) GetAll() ([]*models.Post, error) {
	posts := []*models.Post{}
	query := `SELECT *, (
		SELECT name FROM users WHERE users.id = posts.user_id
		)
	FROM posts`
	rows, err := r.DB.Query(query)
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

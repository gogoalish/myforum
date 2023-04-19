package repository

import (
	"database/sql"
	"errors"

	"forum/internal/models"
)

type Users interface {
	SignUp(models.User) error
	UserByEmail(email string) (models.User, error)
	UserByToken(token string) (models.User, error)
	SetToken(id int, token string) error
}

type UserRepo struct {
	*sql.DB
}

func NewUserRepo(db *sql.DB) Users {
	return &UserRepo{db}
}

func (u *UserRepo) SignUp(m models.User) error {
	query := `INSERT INTO users (email, name, password)
	VALUES(?, ?, ?);`
	if _, err := u.Exec(query, m.Email, m.Name, m.Password); err != nil {
		return err
	}
	return nil
}

func (u *UserRepo) UserByEmail(email string) (models.User, error) {
	query := `SELECT * FROM users
	WHERE ? = email`
	var user models.User
	err := u.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Name, &user.Password, &user.Token)
	if errors.Is(err, sql.ErrNoRows) {
		return user, models.ErrNoRecord
	}
	return user, err
}

func (u *UserRepo) SetToken(id int, token string) error {
	query := `UPDATE users
	SET token = ?
	WHERE ? = id`
	if _, err := u.Exec(query, token, id); err != nil {
		return err
	}
	return nil
}

func (u *UserRepo) UserByToken(token string) (models.User, error) {
	query := `SELECT * FROM users
	WHERE ? = token`
	var user models.User
	err := u.QueryRow(query, token).Scan(&user.ID, &user.Email, &user.Name, &user.Password, &user.Token)
	if errors.Is(err, sql.ErrNoRows) {
		return user, models.ErrNoRecord
	}
	return user, err
}

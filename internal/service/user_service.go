package service

import (
	"errors"
	"fmt"

	"forum/internal/hasher"
	"forum/internal/models"
	"forum/internal/repository"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Users interface {
	SignUp(models.User) error
	SignIn(email, password string) (models.User, error)
	UserByToken(token string) (models.User, error)
	LogOut(token string) error
}

type UserService struct {
	repo repository.Users
}

func (u *UserService) SignUp(m models.User) error {
	user, err := u.repo.UserByEmail(m.Email)
	if err != nil && !errors.Is(err, models.ErrNoRecord) {
		return fmt.Errorf("userservice #1: %w", err)
	}
	if user.Email == m.Email {
		return models.ErrDuplicateEmail
	}
	user, err = u.repo.UserByName(m.Name)
	if err != nil && !errors.Is(err, models.ErrNoRecord) {
		return fmt.Errorf("userservice #2: %w", err)
	}
	if user.Name == m.Name {
		return models.ErrDuplicateName
	}
	m.Password, err = hasher.Encrypt(m.Password)
	if err != nil {
		return fmt.Errorf("userservice #3: %w", err)
	}
	u.repo.SignUp(m)
	return nil
}

func (u *UserService) SignIn(email, password string) (models.User, error) {
	m, err := u.repo.UserByEmail(email)
	if err != nil {
		return m, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(m.Password), []byte(password))
	if err != nil {
		return m, err
	}
	token, err := uuid.NewV4()
	if err != nil {
		return m, err
	}
	tokenstring := token.String()
	m.Token = &tokenstring
	err = u.repo.SetToken(m.ID, *m.Token)
	return m, err
}

func (u *UserService) UserByToken(token string) (models.User, error) {
	user, err := u.repo.UserByToken(token)
	return user, err
}

func (u *UserService) LogOut(token string) error {
	err := u.repo.RemoveToken(token)
	return err
}

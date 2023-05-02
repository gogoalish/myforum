package service

import (
	"errors"
	"fmt"

	"forum/internal/hasher"
	"forum/internal/models"
	"forum/internal/repository"
)

type Users interface {
	SignUp(models.User) error
	SignIn(login, password string) (models.User, error)
	UserByToken(token string) (models.User, error)
	LogOut(token string) error
}

type UserService struct {
	repo repository.Users
}

func (u *UserService) SignUp(m models.User) error {
	user, err := u.repo.UserByEmail(m.Email)
	if err != nil && !errors.Is(err, models.ErrNoRecord) {
		return fmt.Errorf("userservice - signup #1: %w", err)
	}
	if user.Email == m.Email {
		return models.ErrDuplicateEmail
	}
	user, err = u.repo.UserByName(m.Name)
	if err != nil && !errors.Is(err, models.ErrNoRecord) {
		return fmt.Errorf("userservice - signup #2: %w", err)
	}
	if user.Name == m.Name {
		return models.ErrDuplicateName
	}
	m.Password, err = hasher.Encrypt(m.Password)
	if err != nil {
		return fmt.Errorf("userservice - signup #3: %w", err)
	}
	u.repo.SignUp(m)
	return nil
}

func (u *UserService) SignIn(login, password string) (models.User, error) {
	m, err := u.repo.UserByEmail(login)
	switch {
	case errors.Is(err, models.ErrNoRecord):
		m, err = u.repo.UserByName(login)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				return m, err
			}
			return m, fmt.Errorf("userservice - signin #1: %w", err)
		}
	case err != nil:
		return m, fmt.Errorf("userservice - signin #2: %w", err)
	}

	if !hasher.CorrectPassword(m.Password, password) {
		return m, models.ErrInvalidCredentials
	}

	m.Token, err = hasher.GenerateToken()
	if err != nil {
		return m, fmt.Errorf("userservice - signin #3: %w", err)
	}
	err = u.repo.SetToken(m.ID, *m.Token)
	if err != nil {
		return m, fmt.Errorf("userservice - signin #4: %w", err)
	}
	return m, nil
}

func (u *UserService) UserByToken(token string) (models.User, error) {
	return u.repo.UserByToken(token)
}

func (u *UserService) LogOut(token string) error {
	return u.repo.RemoveToken(token)
}

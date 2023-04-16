package service

import (
	"forum/internal/models"
	"forum/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	Users
}

func New(r *repository.Repos) *Service {
	return &Service{
		Users: NewUserService(r),
	}
}

func PasswordCrypt(m *models.User) {
	password := []byte(m.Password)
	crypted, _ := bcrypt.GenerateFromPassword(password, 3)
	m.Password = string(crypted)
}

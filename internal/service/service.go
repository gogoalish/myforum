package service

import (
	"forum/internal/models"
	"forum/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	Users
	Posts
}

func New(r *repository.Repo) *Service {
	return &Service{
		Users: NewUserService(r),
		Posts: &PostService{r.Posts},
	}
}

func Encrypt(m *models.User) {
	password := []byte(m.Password)
	crypted, _ := bcrypt.GenerateFromPassword(password, 3)
	m.Password = string(crypted)
}

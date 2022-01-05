package service

import (
	"log"

	"github.com/sing3demons/go-todos/model"
	"github.com/sing3demons/go-todos/repository"
)

type UserService interface {
	Register(user model.User) error
}

type userService struct {
	Repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		Repo: repo,
	}
}

func (service *userService) Register(user model.User) error {
	err := service.Repo.Register(user)
	if err != nil {
		log.Printf("service todo FindAll, error: %v", err)
		return err
	}
	return nil
}

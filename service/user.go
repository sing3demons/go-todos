package service

import (
	"log"

	"github.com/sing3demons/go-todos/model"
	"github.com/sing3demons/go-todos/repository"
)

type UserService interface {
	Register(user model.User) error
	FindByEmail(email string) (*model.User, error)
	FindByUsers() ([]model.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (service *userService) Register(user model.User) error {
	err := service.repo.Register(user)
	if err != nil {
		log.Printf("service todo FindAll, error: %v", err)
		return err
	}
	return nil
}

func (service *userService) FindByEmail(email string) (*model.User, error) {
	user, err := service.repo.FindByEmail(email)
	if err != nil {
		log.Printf("service todo FindAll, error: %v", err)

		return nil, err
	}
	return user, nil
}

func (service *userService) FindByUsers() ([]model.User, error) {
	user, err := service.repo.FindByUsers()
	if err != nil {
		log.Printf("service user FindByUser, error: %v", err)

		return nil, err
	}
	return user, nil
}

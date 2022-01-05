package repository

import (
	"github.com/sing3demons/go-todos/model"
	"gorm.io/gorm"
)

type UserRepository interface{
	Register(user model.User) error
}

type userRepository struct {
	*gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{DB: db}
}

func (repo *userRepository) Register(user model.User) error {
	err := repo.DB.Create(&user).Error
	if err != nil {
		return err
	}
	return nil
}

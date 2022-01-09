package repository

import (
	"github.com/sing3demons/go-todos/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	Register(user model.User) error

	FindByEmail(email string) (*model.User, error)
	FindByUsers() ([]model.User, error)
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

func (repo *userRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	if err := repo.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *userRepository) FindByUsers() ([]model.User, error) {
	var user []model.User
	if err := repo.DB.Find(&user).Error; err != nil {
		return nil, err
	}
	return user, nil

}

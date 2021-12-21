package repository

import (
	"github.com/sing3demons/go-todos/model"
	"gorm.io/gorm"
)

type TodoRepository interface{
	AllTodos() ([]model.Todo, error)
}

type todoRepository struct{ DB *gorm.DB }

func NewTodoRepository(db *gorm.DB) TodoRepository {
	return &todoRepository{DB: db}
}

func (repo *todoRepository) AllTodos() ([]model.Todo, error) {
	var todos []model.Todo

	if err := repo.DB.Find(&todos).Error; err != nil {
		return nil, err
	}

	return todos, nil
}

package repository

import (
	"github.com/sing3demons/go-todos/model"

	"gorm.io/gorm"
)

type TodoRepository interface {
	AllTodos(limit int, page int) ([]model.Todo, *model.Pagination, error)
	FindTodoByID(id uint) (*model.Todo, error)
	InsertTodo(todo model.Todo) error
	DeleteTodo(todo model.Todo) error
	UpdateTodo(todo model.Todo) error
}

type todoRepository struct{ DB *gorm.DB }

func NewTodoRepository(db *gorm.DB) TodoRepository {
	return &todoRepository{DB: db}
}

func (repo *todoRepository) AllTodos(limit int, page int) ([]model.Todo, *model.Pagination, error) {
	var todos []model.Todo

	pagination := model.Pagination{
		Limit: limit,
		Page:  page,
	}

	if err := repo.DB.Scopes(repo.paginate(&todos, &pagination)).Find(&todos).Error; err != nil {
		return nil, nil, err
	}

	return todos, &pagination, nil
}

func (repo *todoRepository) FindTodoByID(id uint) (*model.Todo, error) {
	var todo model.Todo

	if err := repo.DB.First(&todo, id).Error; err != nil {
		return nil, err
	}
	return &todo, nil
}

func (repo *todoRepository) InsertTodo(todo model.Todo) error {
	if err := repo.DB.Create(&todo).Error; err != nil {
		return err
	}

	return nil
}

func (repo *todoRepository) DeleteTodo(todo model.Todo) error {
	if err := repo.DB.Delete(&todo).Error; err != nil {
		return err
	}
	return nil
}

func (repo *todoRepository) UpdateTodo(todo model.Todo) error {
	if err := repo.DB.Save(&todo).Error; err != nil {
		return err
	}
	return nil
}

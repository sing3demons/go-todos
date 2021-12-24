package repository

import (
	"math"

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

	if err := repo.DB.Scopes(paginate(&todos, &pagination, repo.DB)).Find(&todos).Error; err != nil {
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

func paginate(value interface{}, pagination *model.Pagination, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	var totalRows int64
	db.Model(value).Count(&totalRows)
	pagination.TotalRows = totalRows
	totalPages := int(math.Ceil(float64(totalRows) / float64(pagination.Limit)))
	pagination.TotalPages = totalPages

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.Limit).Order(pagination.GetSort())
	}
}

package service

import (
	"fmt"

	"github.com/sing3demons/go-todos/model"
	"github.com/sing3demons/go-todos/pagination"
	"github.com/sing3demons/go-todos/repository"
)

type TodoService interface {
	FindTodos(limit int, page int) ([]model.Todo, *pagination.Pagination, error)
	FindTodo(id uint) (*model.Todo, error)
	CreateTodo(todo model.Todo) error
	DeleteTodo(todo model.Todo) error
	UpdateTodo(todo model.Todo) error
}

type todoService struct {
	repo repository.TodoRepository
}

func NewTodoService(repo repository.TodoRepository) TodoService {
	return &todoService{repo: repo}
}

func (s *todoService) FindTodos(limit int, page int) ([]model.Todo, *pagination.Pagination, error) {
	todos, paging, err := s.repo.AllTodos(limit, page)
	if err != nil {
		fmt.Printf("error: %v", err)
		return nil, nil, err
	}
	return todos, paging, nil
}

func (s *todoService) FindTodo(id uint) (*model.Todo, error) {
	todo, err := s.repo.FindTodoByID(id)
	if err != nil {
		fmt.Printf("error: %v", err)
		return nil, err
	}
	return todo, nil
}

func (s *todoService) CreateTodo(todo model.Todo) error {
	err := s.repo.InsertTodo(todo)
	if err != nil {
		fmt.Printf("error: %v", err)
		return err
	}
	return nil
}

func (s *todoService) DeleteTodo(todo model.Todo) error {
	err := s.repo.DeleteTodo(todo)
	if err != nil {
		fmt.Printf("error: %v", err)
		return err
	}
	return nil
}

func (s *todoService) UpdateTodo(todo model.Todo) error {
	err := s.repo.UpdateTodo(todo)
	if err != nil {
		fmt.Printf("error: %v", err)
		return err
	}
	return nil
}

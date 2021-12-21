package service

import (
	"fmt"

	"github.com/sing3demons/go-todos/model"
	"github.com/sing3demons/go-todos/repository"
)

type TodoService interface{
	FindTodos() ([]model.Todo, error)
}

type todoService struct {
	repo repository.TodoRepository
}

func NewTodoService(repo repository.TodoRepository) TodoService {
	return &todoService{repo: repo}
}

func (s *todoService) FindTodos() ([]model.Todo, error) {
	todos, err := s.repo.AllTodos()
	if err != nil {
		fmt.Printf("error: %v", err)
		return nil, err
	}
	return todos, nil
}

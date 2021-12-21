package repository

type TodoRepository interface{}

type todoRepository struct{}

func NewTodoRepository() TodoRepository {
	return &todoRepository{}
}

func (repo *todoRepository) AllTodos() {}

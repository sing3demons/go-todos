package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sing3demons/go-todos/handler"
	"github.com/sing3demons/go-todos/redis"
	"github.com/sing3demons/go-todos/repository"
	"github.com/sing3demons/go-todos/service"
	"gorm.io/gorm"
)

type Router struct {
	*fiber.App
	*gorm.DB
	*redis.Cacher
}

func (r *Router) Serve() {
	v1 := r.App.Group("/api/v1/")

	todoGroup := v1.Group("todos")
	r.todoRouter(todoGroup)

}

func (r *Router) todoRouter(todoGroup fiber.Router) {
	todoRepository := repository.NewTodoRepository(r.DB)
	todoService := service.NewTodoService(todoRepository)
	todoHandler := handler.NewtodoHandler(todoService, r.Cacher)

	todoGroup.Get("", todoHandler.AllTodos)
	todoGroup.Get("/:id", todoHandler.FindTodo)
	todoGroup.Post("", todoHandler.CreateTodo)
	todoGroup.Delete("/:id", todoHandler.DeleteTodo)
	todoGroup.Put("/:id", todoHandler.UpdateTodo)
}

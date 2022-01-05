package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sing3demons/go-todos/cache"
	"github.com/sing3demons/go-todos/handler"
	"github.com/sing3demons/go-todos/repository"
	"github.com/sing3demons/go-todos/service"
	"gorm.io/gorm"
)

type Router struct {
	*fiber.App
	*gorm.DB
	*cache.Cacher
}

func (r *Router) Serve() {
	v1 := r.App.Group("/api/v1/")

	todoGroup := v1.Group("todos")
	r.todoRouter(todoGroup)

	authGroup := v1.Group("auth")
	r.authRouter(authGroup)

}

func (r *Router) authRouter(authGroup fiber.Router) {
	repository := repository.NewUserRepository(r.DB)
	service := service.NewUserService(repository)
	handler := handler.NewUserHandler(service)

	authGroup.Post("/sign-up", handler.Register)
}

func (r *Router) todoRouter(todoGroup fiber.Router) {
	repository := repository.NewTodoRepository(r.DB)
	todoService := service.NewTodoService(repository)
	todoHandler := handler.NewtodoHandler(todoService, r.Cacher)

	todoGroup.Get("", todoHandler.AllTodos)
	todoGroup.Get("/:id", todoHandler.FindTodo)
	todoGroup.Post("", todoHandler.CreateTodo)
	todoGroup.Delete("/:id", todoHandler.DeleteTodo)
	todoGroup.Put("/:id", todoHandler.UpdateTodo)
}

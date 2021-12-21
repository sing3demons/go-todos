package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sing3demons/go-todos/database"
	"github.com/sing3demons/go-todos/handler"
	"github.com/sing3demons/go-todos/repository"
	"github.com/sing3demons/go-todos/service"
)

func Serve(app *fiber.App) {
	v1 := app.Group("/api/v1/")

	db := database.GetDB()

	todoRepository := repository.NewTodoRepository(db)
	todoService := service.NewTodoService(todoRepository)
	todoHandler := handler.NewtodoHandler(todoService)

	todoGroup := v1.Group("todos")
	todoGroup.Get("", todoHandler.AllTodos)
}

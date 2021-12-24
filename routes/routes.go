package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sing3demons/go-todos/database"
	"github.com/sing3demons/go-todos/handler"
	"github.com/sing3demons/go-todos/repository"
	"github.com/sing3demons/go-todos/service"
	"gorm.io/gorm"
)

func Serve(app *fiber.App) {
	v1 := app.Group("/api/v1/")

	db := database.GetDB()

	todoGroup := v1.Group("todos")
	todoRouter(todoGroup, db)

}

func todoRouter(todoGroup fiber.Router, db *gorm.DB) {
	todoRepository := repository.NewTodoRepository(db)
	todoService := service.NewTodoService(todoRepository)
	todoHandler := handler.NewtodoHandler(todoService)

	todoGroup.Get("", todoHandler.AllTodos)
	todoGroup.Get("/:id", todoHandler.FindTodo)
	todoGroup.Post("", todoHandler.CreateTodo)
	todoGroup.Delete("/:id", todoHandler.DeleteTodo)
	todoGroup.Put("/:id", todoHandler.UpdateTodo)
}

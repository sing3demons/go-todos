package routes

import (
	"github.com/sing3demons/go-todos/cache"
	"github.com/sing3demons/go-todos/database"
	"github.com/sing3demons/go-todos/handler"
	"github.com/sing3demons/go-todos/middleware"
	"github.com/sing3demons/go-todos/repository"
	"github.com/sing3demons/go-todos/service"
)

func Serve(app *MyRouter) {
	v1 := app.Group("/api/v1/")
	db := database.GetDB()
	cache := cache.NewCacher(&cache.CacherConfig{})

	todoGroup := v1.Group("todos")
	{
		repository := repository.NewTodoRepository(db)
		todoService := service.NewTodoService(repository)
		todoHandler := handler.NewTodoHandler(todoService, cache)

		todoGroup.Get("", todoHandler.AllTodos)
		todoGroup.Get("/:id", todoHandler.FindTodo)
		todoGroup.Post("", todoHandler.CreateTodo)
		todoGroup.Delete("/:id", todoHandler.DeleteTodo)
		todoGroup.Put("/:id", todoHandler.UpdateTodo)
	}

	authGroup := v1.Group("auth")
	{
		repository := repository.NewUserRepository(db)
		service := service.NewUserService(repository)
		handler := handler.NewUserHandler(service)

		authenticate := middleware.JwtVerify()

		authGroup.Post("/sign-up", handler.Register)
		authGroup.Post("/sign-in", handler.Login)
		authGroup.Use(authenticate)
		authGroup.Get("/profile", handler.Profile)
	}

	userGroup := v1.Group("users")
	{
		repository := repository.NewUserRepository(db)
		service := service.NewUserService(repository)
		handler := handler.NewUserHandler(service)

		authenticate := middleware.JwtVerify()
		authorize := middleware.Authorize()
		userGroup.Get("/", authenticate, authorize, handler.FindUsers)
	}

}

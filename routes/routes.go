package routes

import "github.com/gofiber/fiber/v2"

func Serve(app *fiber.App) {
	v1 := app.Group("/api/v1/")

	todoGroup := v1.Group("todos")
	todoGroup.Get("", func(c *fiber.Ctx) error {
		return c.Status(200).JSON("hello, world")
	})
}
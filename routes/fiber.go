package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type MyRouter struct {
	*fiber.App
}

func NewFiberRouter() *MyRouter {
	app := fiber.New()
	app.Use(recover.New())

	// app.Use(cors.New(cors.Config{
	// 	AllowOrigins: "http://localhost:8080, http://127.0.0.1:8080/swagger/index.html",
	// 	AllowHeaders: "Origin, Content-Type, Accept",
	// }))

	app.Use(cors.New())

	return &MyRouter{app}
}

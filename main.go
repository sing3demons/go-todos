package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"

	"github.com/sing3demons/go-todos/database"
	"github.com/sing3demons/go-todos/redis"
	"github.com/sing3demons/go-todos/routes"
)

var (
	buildcommit = "dev"
	buildtime   = time.Now().String()
)

func NewCacherConfig() *redis.CacherConfig {
	return &redis.CacherConfig{}
}

func main() {
	// Liveness Probe
	_, err := os.Create("/tmp/live")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove("/tmp/live")

	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load(".env")
		if err != nil {
			log.Println("Error loading .env file")
		}
	}
	// connect database
	database.InitDB()
	// seeds.Load()

	app := fiber.New()
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:8080",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	app.Get("/dashboard", monitor.New())
	app.Static("/uploads", "./uploads")

	//สร้าง folder
	uploadDirs := [...]string{"todo", "users"}
	for _, dir := range uploadDirs {
		path := fmt.Sprintf("uploads/%s/images", dir)
		os.MkdirAll(path, 0755)
	}

	if os.Getenv("APP_ENV") == "production" {
		file, err := os.OpenFile("./logs/logs.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		defer file.Close()

		app.Use(logger.New(logger.Config{
			Output:       file,
			Format:       "[${time}], ${status} - ${latency}, ip:(${ip}:${pid}), ${method}, ${path}\n",
			Next:         nil,
			TimeFormat:   "15:04:05",
			TimeZone:     "Local",
			TimeInterval: 500 * time.Millisecond,
		}))
	}

	app.Use(logger.New(logger.ConfigDefault))
	app.Get("/x", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"buildcommit": buildcommit,
			"buildtime":   buildtime,
		})
	})
	// Readiness Probe
	app.Get("/healthz", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })

	//Router
	db := database.GetDB()
	redis := redis.NewCacher(NewCacherConfig())
	r := routes.Router{
		App:   app,
		DB:    db,
		Redis: redis,
	}
	r.Serve()

	//Graceful Shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := app.Listen(":" + os.Getenv("PORT")); err != nil {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-ctx.Done()
	stop()

	fmt.Println("shutting down gracefully, press Ctrl+C again to force")

	if err := app.Shutdown(); err != nil {
		fmt.Println(err)
	}
}

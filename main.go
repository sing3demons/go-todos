package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	swagger "github.com/arsmn/fiber-swagger/v2" // fiber-swagger middleware
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/joho/godotenv"

	"github.com/sing3demons/go-todos/database"
	_ "github.com/sing3demons/go-todos/docs"
	"github.com/sing3demons/go-todos/routes"
	"github.com/sing3demons/go-todos/seeds"
)

var (
	buildcommit = "dev"
	buildtime   = time.Now().String()
	production  = "production"
)

func init() {
	if os.Getenv("APP_ENV") != production {
		err := godotenv.Load(".env")
		if err != nil {
			log.Println("Error loading .env file")
		}
	}

	if os.Getenv("APP_ENV") != production {
		database.InitDB()
		seeds.Load()
	}
}

// @title Fiber go-todos API
// @version 1.0
// @description This is a sample swagger for Fiber

// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @host localhost:8080
// @BasePath /
func main() {
	var port string = os.Getenv("PORT")
	// Liveness Probe
	_, err := os.Create("/tmp/live")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove("/tmp/live")

	// connect database
	database.InitDB()

	app := routes.NewFiberRouter()
	app.Get("/swagger/*", swagger.HandlerDefault)

	app.Get("/dashboard", monitor.New())
	app.Static("/uploads", "./uploads")

	//สร้าง folder
	uploadDirs := [...]string{"todo", "users"}
	for _, dir := range uploadDirs {
		path := fmt.Sprintf("uploads/%s/images", dir)
		os.MkdirAll(path, 0755)
	}

	app.Use(logger.New(logger.ConfigDefault))
	app.Get("/x", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"buildcommit": buildcommit,
			"buildtime":   buildtime,
		})
	})
	// Readiness Probe
	app.Get("/healthz", healthz)

	//Router
	routes.Serve(app)

	//Graceful Shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		fmt.Printf("\nBrowse to http://127.0.0.1:%s/swagger/index.html\n", port)

		if err := app.Listen(":" + port); err != nil {
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

func downloadLogFile(c *fiber.Ctx) error {
	url := "./logs/logs.log"

	// resp, err := os.Open(url)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// defer resp.Close()

	// body, err := ioutil.ReadAll(resp)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// data := []byte("")

	// fmt.Println(string(body))
	// os.WriteFile(url, data, 0666)

	return c.Download(url, "log")
}

// ShowHealth godoc
// @Summary Show a healthz
// @Description get healthz
// @Accept  json
// @Produce  json
// @Success 200
// @Router /healthz [get]
func healthz(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) }

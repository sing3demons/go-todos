package seeds

import (
	"fmt"
	"os"
	"strconv"

	"github.com/bxcodec/faker/v3"
	"github.com/sing3demons/go-todos/database"
	"github.com/sing3demons/go-todos/model"
)

func Load() {
	db := database.GetDB()
	numOfTodo := 5000
	// db.Migrator().DropTable(&model.Todo{})
	// db.AutoMigrate(&model.Todo{})

	todos := make([]model.Todo, numOfTodo)
	var count int64

	if os.Getenv("APP_ENV") == "dev" {
		db.Find(&todos).Count(&count)
		if count > 0 {
			return
		}

		fmt.Println("seed data")

		for i := 1; i <= numOfTodo; i++ {
			todo := model.Todo{
				Title: faker.Name(),
				Desc:  faker.Word(),
				Image: "https://source.unsplash.com/random/300x200?" + strconv.Itoa(i),
			}
			todos = append(todos, todo)
		}
		db.Create(&todos)
		fmt.Println("end...")
	}

	if os.Getenv("APP_ENV") == "dev" {
		var user model.User
		db.Find(&user).Count(&count)
		if count > 0 {
			return
		}

		fmt.Println("seed user")

		user = model.User{
			Email:     "admin@dev.com",
			Password:  "p@ssw0rd",
			FirstName: "admin",
			Role:      "Admin",
		}
		db.Create(&user)
		fmt.Println("end...")
	}
}

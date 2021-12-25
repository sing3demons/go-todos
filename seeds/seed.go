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

	if os.Getenv("APP_ENV") == "dev" {
		fmt.Println("start")
		numOfTodo := 50000
		// db.Migrator().DropTable(&model.Todo{})
		// db.AutoMigrate(&model.Todo{})

		todos := make([]model.Todo, numOfTodo)

		for i := 1; i <= numOfTodo; i++ {
			todo := model.Todo{
				Title: faker.Name(),
				Desc:  faker.Word(),
				Image: "https://source.unsplash.com/random/300x200?" + strconv.Itoa(i),
			}
			todos = append(todos, todo)
		}
		db.CreateInBatches(todos, 100)
		fmt.Println("end...")
	}
}

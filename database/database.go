package database

import (
	"os"

	"github.com/sing3demons/go-todos/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB() {
	// user := os.Getenv("DB_USER")
	// password := os.Getenv("DB_PASSWORD")
	// dbName := os.Getenv("DB_NAME")
	// port := os.Getenv("DB_PORT")
	// host := os.Getenv("DB_HOST")

	// dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s  sslmode=disable TimeZone=Asia/Bangkok", host, user, password, dbName, port)
	// dsn := fmt.Sprint("host=localhost user=postgresadmin password=admin123 dbname=postgresdb port=30001  sslmode=disable TimeZone=Asia/Bangkok")
	// dsn := "host=db user=postgres password=passw0rd dbname=todos port=5432  sslmode=disable TimeZone=Asia/Bangkok"
	dsn := os.Getenv("DSN")
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// database.Migrator().DropTable(&model.Todo{})
	database.AutoMigrate(&model.Todo{})
	database.AutoMigrate(&model.User{})

	db = database
}

func GetDB() *gorm.DB {
	return db
}

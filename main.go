package main

import (
	"fmt"
	"log"
	"restapi/book"
	"restapi/database"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func initDatabase() {
	var err error
	database.DBConn, err = gorm.Open(sqlite.Open("books.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	fmt.Println("DB connection successful")

	if err := database.DBConn.AutoMigrate(&book.Book{}); err != nil {
		log.Fatal("Failed to migrate DB:", err)
	}
	fmt.Println("Database migrated")
}

func main() {
	// Fiber instance
	app := fiber.New()

	// Init DB
	initDatabase()
	defer func() {
		sqlDB, err := database.DBConn.DB()
		if err == nil {
			sqlDB.Close()
		}
	}()

	// Routes
	handleRoute(app)

	// Start server
	log.Fatal(app.Listen(":3000"))
}

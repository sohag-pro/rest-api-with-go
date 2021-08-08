package main

import (
	"fmt"
	"restapi/database"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)



func initDatabase(){
	var err error
	database.DBConn, err = gorm.Open("sqlite3", "books.db")

	if err != nil {
		panic("Failed to connect to DB")
	}

	fmt.Println("DB connection Successfull")
}

func main(){
	// Fiber Instance
	app := fiber.New()

	// Init DB
	initDatabase()
	defer database.DBConn.Close()

	// Route
	handleRoute(app)


	// Start Server
	app.Listen(":3000")
}



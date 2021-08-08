package main

import "github.com/gofiber/fiber/v2"

func main(){
	// Fiber Instance
	app := fiber.New()

	// Route
	handleRoute(app)


	// Start Server
	app.Listen(":3000")
}



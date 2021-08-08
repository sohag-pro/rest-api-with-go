package main

import (
	"restapi/book"

	"github.com/gofiber/fiber/v2"
)

func handleRoute(app *fiber.App){
	app.Get("api/v1/book", book.GetBooks)
	app.Get("api/v1/book/:id", book.GetBook)
	app.Post("api/v1/book", book.NewBook)
	app.Patch("api/v1/book/:id", book.UpdateBook)
	app.Delete("api/v1/book/:id", book.DeleteBooks)
}

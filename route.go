package main

import (
	"restapi/book"

	"github.com/gofiber/fiber/v2"
)

func handleRoute(app *fiber.App, cfg Config) {
	api := app.Group("/api/v1")

	// Public read endpoints
	api.Get("/book", book.GetBooks)
	api.Get("/book/:id", book.GetBook)

	// Mutating endpoints, guarded by API key when configured
	write := api.Group("", apiKeyAuth(cfg.APIKey))
	write.Post("/book", book.NewBook)
	write.Patch("/book/:id", book.UpdateBook)
	write.Delete("/book/:id", book.DeleteBooks)
}

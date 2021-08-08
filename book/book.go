package book

import "github.com/gofiber/fiber/v2"

func GetBooks(c *fiber.Ctx) error {
	return c.SendString("All Books")
}

func GetBook(c *fiber.Ctx) error {
	return c.SendString("A single Book")
}

func NewBook(c *fiber.Ctx) error {
	return c.SendString("New Book")
}

func DeleteBooks(c *fiber.Ctx) error {
	return c.SendString("Delete Book")
}


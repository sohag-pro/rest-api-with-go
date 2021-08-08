package book

import (
	"restapi/database"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/gorm"
)


// Book Model
type Book struct {
	gorm.Model
	Title string `json:"title"`
	Author string `json:"author"`
	Rating int `json:"rating"`
}


// Get all books list
func GetBooks(c *fiber.Ctx) error {
	db := database.DBConn
	var books []Book

	db.Find(&books)

	return c.JSON(books)
}


// Get a single book details
func GetBook(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn
	var book Book

	db.Find(&book, id)

	if book.Title == ""{
		return c.Status(404).SendString("Book not found")
	}

	return c.JSON(book)
}


// Create a new book
func NewBook(c *fiber.Ctx) error {
	db := database.DBConn
	book := new(Book)

	if err := c.BodyParser(book); err != nil {
		return c.Status(406).Send([]byte(err.Error()))
	} 

	if book.Title == ""{
		return c.Status(400).SendString("Title is required")
	}

	db.Create(&book)

	return c.JSON(book)
}


// Update a book details
func UpdateBook(c *fiber.Ctx) error {
	db := database.DBConn
	id := c.Params("id")
	updated_book := new(Book)
	var book Book

	db.First(&book, id)

	if err := c.BodyParser(updated_book); err != nil {
		return c.Status(406).Send([]byte(err.Error()))
	} 

	book.Title = updated_book.Title
	book.Author = updated_book.Author
	book.Rating = updated_book.Rating

	db.Save(&book)

	return c.JSON(book)
}


// Delete a book
func DeleteBooks(c *fiber.Ctx) error {
	db := database.DBConn
	id := c.Params("id")
	var book Book

	db.First(&book, id)

	if book.Title == ""{
		return c.Status(404).SendString("Book not found")
	}

	db.Delete(&book)

	return c.SendString("Book Deleted Successfully")
}


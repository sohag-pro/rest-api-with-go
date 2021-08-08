package book

import (
	"restapi/database"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/gorm"
)


type Book struct {
	gorm.Model
	Title string `json:"title"`
	Author string `json:"author"`
	Rating int `json:"rating"`
}

func GetBooks(c *fiber.Ctx) error {
	db := database.DBConn
	var books []Book

	db.Find(&books)

	return c.JSON(books)
}

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

func NewBook(c *fiber.Ctx) error {
	db := database.DBConn
	book := new(Book)

	if err := c.BodyParser(book); err != nil {
		return c.Status(406).Send([]byte(err.Error()))
	} 

	db.Create(&book)

	return c.JSON(book)
}

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


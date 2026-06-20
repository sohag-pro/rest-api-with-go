package book

import (
	"errors"
	"strings"

	"restapi/database"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Book Model
type Book struct {
	gorm.Model
	Title  string `json:"title"`
	Author string `json:"author"`
	Rating int    `json:"rating"`
}

// validate normalizes a book and returns an error message if invalid.
func (b *Book) validate() string {
	b.Title = strings.TrimSpace(b.Title)
	b.Author = strings.TrimSpace(b.Author)

	if b.Title == "" {
		return "title is required"
	}
	if b.Rating < 0 || b.Rating > 5 {
		return "rating must be between 0 and 5"
	}
	return ""
}

// findByID loads a book by id; returns false if it does not exist.
func findByID(db *gorm.DB, id string, book *Book) bool {
	err := db.First(book, "id = ?", id).Error
	return !errors.Is(err, gorm.ErrRecordNotFound)
}

// GetBooks lists books with optional pagination: ?limit=&offset=
// limit defaults to 10 (max 100), offset defaults to 0.
func GetBooks(c *fiber.Ctx) error {
	db := database.DBConn
	var books []Book

	limit := c.QueryInt("limit", 10)
	if limit < 1 || limit > 100 {
		limit = 10
	}
	offset := c.QueryInt("offset", 0)
	if offset < 0 {
		offset = 0
	}

	db.Limit(limit).Offset(offset).Find(&books)

	return c.JSON(books)
}

// GetBook returns a single book by id.
func GetBook(c *fiber.Ctx) error {
	db := database.DBConn
	var book Book

	if !findByID(db, c.Params("id"), &book) {
		return c.Status(404).SendString("Book not found")
	}

	return c.JSON(book)
}

// NewBook creates a book.
func NewBook(c *fiber.Ctx) error {
	db := database.DBConn
	book := new(Book)

	if err := c.BodyParser(book); err != nil {
		return c.Status(406).Send([]byte(err.Error()))
	}

	if msg := book.validate(); msg != "" {
		return c.Status(400).SendString(msg)
	}

	db.Create(&book)

	return c.JSON(book)
}

// UpdateBook updates an existing book.
func UpdateBook(c *fiber.Ctx) error {
	db := database.DBConn
	var book Book

	if !findByID(db, c.Params("id"), &book) {
		return c.Status(404).SendString("Book not found")
	}

	updated := new(Book)
	if err := c.BodyParser(updated); err != nil {
		return c.Status(406).Send([]byte(err.Error()))
	}

	book.Title = updated.Title
	book.Author = updated.Author
	book.Rating = updated.Rating

	if msg := book.validate(); msg != "" {
		return c.Status(400).SendString(msg)
	}

	db.Save(&book)

	return c.JSON(book)
}

// DeleteBooks deletes a book by id.
func DeleteBooks(c *fiber.Ctx) error {
	db := database.DBConn
	var book Book

	if !findByID(db, c.Params("id"), &book) {
		return c.Status(404).SendString("Book not found")
	}

	db.Delete(&book)

	return c.SendString("Book Deleted Successfully")
}

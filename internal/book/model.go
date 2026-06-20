// Package book implements the book domain: model, repository, service, and HTTP handlers.
package book

import (
	"errors"
	"strings"

	"gorm.io/gorm"
)

// ErrNotFound is returned when a book does not exist.
var ErrNotFound = errors.New("book not found")

// ValidationError indicates invalid book input. It maps to HTTP 400.
type ValidationError struct {
	Msg string
}

func (e ValidationError) Error() string { return e.Msg }

// Book is the domain model and database row.
type Book struct {
	gorm.Model
	Title  string `json:"title"`
	Author string `json:"author"`
	Rating int    `json:"rating"`
}

// Normalize trims user-supplied string fields in place.
func (b *Book) Normalize() {
	b.Title = strings.TrimSpace(b.Title)
	b.Author = strings.TrimSpace(b.Author)
}

// Validate returns a ValidationError if the book is invalid.
func (b *Book) Validate() error {
	if b.Title == "" {
		return ValidationError{Msg: "title is required"}
	}
	if b.Rating < 0 || b.Rating > 5 {
		return ValidationError{Msg: "rating must be between 0 and 5"}
	}
	return nil
}

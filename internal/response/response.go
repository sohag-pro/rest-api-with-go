// Package response provides a consistent JSON envelope for API errors.
package response

import "github.com/gofiber/fiber/v2"

// ErrorBody is the standard error payload returned to clients.
type ErrorBody struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

// Error writes a JSON error response with the given HTTP status.
func Error(c *fiber.Ctx, status int, msg string) error {
	return c.Status(status).JSON(ErrorBody{Error: msg, Code: status})
}

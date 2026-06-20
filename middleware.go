package main

import "github.com/gofiber/fiber/v2"

// apiKeyAuth guards mutating routes with a shared API key.
// When key is empty, auth is disabled and all requests pass through.
// Clients authenticate via the "X-API-Key" request header.
func apiKeyAuth(key string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if key == "" {
			return c.Next()
		}
		if c.Get("X-API-Key") != key {
			return c.Status(401).SendString("Unauthorized")
		}
		return c.Next()
	}
}

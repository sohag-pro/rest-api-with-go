package server

import (
	"log/slog"
	"time"

	"restapi/internal/response"

	"github.com/gofiber/fiber/v2"
)

// apiKeyAuth guards mutating routes with a shared API key.
// When key is empty, auth is disabled. Clients send the key in X-API-Key.
func apiKeyAuth(key string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if key == "" {
			return c.Next()
		}
		if c.Get("X-API-Key") != key {
			return response.Error(c, fiber.StatusUnauthorized, "unauthorized")
		}
		return c.Next()
	}
}

// requestLogger emits a structured slog line per request.
func requestLogger(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()

		// On unmatched routes the framework error handler runs after this
		// middleware, so derive the real status from the returned error.
		status := c.Response().StatusCode()
		if fe, ok := err.(*fiber.Error); ok {
			status = fe.Code
		}

		logger.Info("request",
			"method", c.Method(),
			"path", c.Path(),
			"status", status,
			"duration_ms", time.Since(start).Milliseconds(),
			"request_id", c.GetRespHeader(fiber.HeaderXRequestID),
			"ip", c.IP(),
		)
		return err
	}
}

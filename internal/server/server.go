// Package server builds the configured Fiber application.
package server

import (
	"log/slog"

	"restapi/internal/book"
	"restapi/internal/config"
	"restapi/internal/response"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"gorm.io/gorm"
)

// New constructs a Fiber app with middleware, health checks, and routes wired.
func New(cfg config.Config, db *gorm.DB, logger *slog.Logger) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:               "books-api",
		ReadTimeout:           cfg.ReadTimeout,
		WriteTimeout:          cfg.WriteTimeout,
		DisableStartupMessage: true,
		ErrorHandler:          errorHandler,
	})

	app.Use(requestid.New())
	app.Use(recover.New())
	app.Use(cors.New())
	app.Use(requestLogger(logger))

	// Liveness/readiness probe.
	app.Get("/healthz", func(c *fiber.Ctx) error {
		sqlDB, err := db.DB()
		if err != nil || sqlDB.Ping() != nil {
			return response.Error(c, fiber.StatusServiceUnavailable, "database unavailable")
		}
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// API v1.
	repo := book.NewRepository(db)
	svc := book.NewService(repo)
	handler := book.NewHandler(svc)
	handler.Register(app.Group("/api/v1"), apiKeyAuth(cfg.APIKey))

	return app
}

// errorHandler is the catch-all for unhandled errors and 404s.
func errorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	msg := "internal server error"
	var fe *fiber.Error
	if e, ok := err.(*fiber.Error); ok {
		fe = e
		code = fe.Code
		msg = fe.Message
	}
	return response.Error(c, code, msg)
}

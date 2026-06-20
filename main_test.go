package main

import (
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func authApp(key string) *fiber.App {
	app := fiber.New()
	app.Use(apiKeyAuth(key))
	app.Get("/x", func(c *fiber.Ctx) error { return c.SendString("ok") })
	return app
}

func TestApiKeyAuthDisabled(t *testing.T) {
	app := authApp("") // empty key disables auth
	req, _ := http.NewRequest("GET", "/x", nil)
	resp, _ := app.Test(req, -1)
	if resp.StatusCode != 200 {
		t.Fatalf("disabled auth want 200, got %d", resp.StatusCode)
	}
}

func TestApiKeyAuthMissingKey(t *testing.T) {
	app := authApp("secret")
	req, _ := http.NewRequest("GET", "/x", nil)
	resp, _ := app.Test(req, -1)
	if resp.StatusCode != 401 {
		t.Fatalf("missing key want 401, got %d", resp.StatusCode)
	}
}

func TestApiKeyAuthWrongKey(t *testing.T) {
	app := authApp("secret")
	req, _ := http.NewRequest("GET", "/x", nil)
	req.Header.Set("X-API-Key", "nope")
	resp, _ := app.Test(req, -1)
	if resp.StatusCode != 401 {
		t.Fatalf("wrong key want 401, got %d", resp.StatusCode)
	}
}

func TestApiKeyAuthValidKey(t *testing.T) {
	app := authApp("secret")
	req, _ := http.NewRequest("GET", "/x", nil)
	req.Header.Set("X-API-Key", "secret")
	resp, _ := app.Test(req, -1)
	if resp.StatusCode != 200 {
		t.Fatalf("valid key want 200, got %d", resp.StatusCode)
	}
}

func TestLoadConfigDefaults(t *testing.T) {
	t.Setenv("PORT", "")
	t.Setenv("DB_PATH", "")
	t.Setenv("API_KEY", "")
	cfg := loadConfig()
	if cfg.Port != "3000" || cfg.DBPath != "books.db" || cfg.APIKey != "" {
		t.Fatalf("unexpected defaults: %+v", cfg)
	}
}

func TestLoadConfigOverrides(t *testing.T) {
	t.Setenv("PORT", "8080")
	t.Setenv("DB_PATH", "/tmp/x.db")
	t.Setenv("API_KEY", "abc")
	cfg := loadConfig()
	if cfg.Port != "8080" || cfg.DBPath != "/tmp/x.db" || cfg.APIKey != "abc" {
		t.Fatalf("overrides not applied: %+v", cfg)
	}
}

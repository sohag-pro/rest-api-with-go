package server

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"path/filepath"
	"testing"
	"time"

	"restapi/internal/config"
	"restapi/internal/database"

	"github.com/gofiber/fiber/v2"
)

const defaultTimeout = 10 * time.Second

func discardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func newTestApp(t *testing.T, apiKey string) *fiber.App {
	t.Helper()
	cfg := config.Config{
		Port:            "3000",
		DBPath:          filepath.Join(t.TempDir(), "test.db"),
		APIKey:          apiKey,
		LogLevel:        "error",
		ReadTimeout:     defaultTimeout,
		WriteTimeout:    defaultTimeout,
		ShutdownTimeout: defaultTimeout,
	}
	db, err := database.Open(cfg)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { _ = database.Close(db) })
	return New(cfg, db, discardLogger())
}

func req(t *testing.T, app *fiber.App, method, url, key string, body any) *http.Response {
	t.Helper()
	var rdr io.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		rdr = bytes.NewReader(b)
	}
	r, _ := http.NewRequest(method, url, rdr)
	if body != nil {
		r.Header.Set("Content-Type", "application/json")
	}
	if key != "" {
		r.Header.Set("X-API-Key", key)
	}
	resp, err := app.Test(r, -1)
	if err != nil {
		t.Fatalf("%s %s: %v", method, url, err)
	}
	return resp
}

func decode[T any](t *testing.T, resp *http.Response) T {
	t.Helper()
	var v T
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return v
}

func TestHealthz(t *testing.T) {
	app := newTestApp(t, "")
	resp := req(t, app, "GET", "/healthz", "", nil)
	if resp.StatusCode != 200 {
		t.Fatalf("healthz want 200, got %d", resp.StatusCode)
	}
}

func TestCRUDFlow(t *testing.T) {
	app := newTestApp(t, "")

	created := req(t, app, "POST", "/api/v1/book", "", map[string]any{"title": "Go", "author": "Pike", "rating": 5})
	if created.StatusCode != 201 {
		t.Fatalf("create want 201, got %d", created.StatusCode)
	}
	b := decode[map[string]any](t, created)
	if b["title"] != "Go" {
		t.Fatalf("bad create body: %v", b)
	}

	got := req(t, app, "GET", "/api/v1/book/1", "", nil)
	if got.StatusCode != 200 {
		t.Fatalf("get want 200, got %d", got.StatusCode)
	}

	upd := req(t, app, "PATCH", "/api/v1/book/1", "", map[string]any{"title": "Go2", "rating": 4})
	if upd.StatusCode != 200 {
		t.Fatalf("update want 200, got %d", upd.StatusCode)
	}

	del := req(t, app, "DELETE", "/api/v1/book/1", "", nil)
	if del.StatusCode != 200 {
		t.Fatalf("delete want 200, got %d", del.StatusCode)
	}

	after := req(t, app, "GET", "/api/v1/book/1", "", nil)
	if after.StatusCode != 404 {
		t.Fatalf("after delete want 404, got %d", after.StatusCode)
	}
}

func TestValidationErrorEnvelope(t *testing.T) {
	app := newTestApp(t, "")
	resp := req(t, app, "POST", "/api/v1/book", "", map[string]any{"title": "", "rating": 9})
	if resp.StatusCode != 400 {
		t.Fatalf("want 400, got %d", resp.StatusCode)
	}
	body := decode[map[string]any](t, resp)
	if body["error"] == nil || body["code"] == nil {
		t.Fatalf("missing error envelope fields: %v", body)
	}
}

func TestUpdateNotFoundNoPhantom(t *testing.T) {
	app := newTestApp(t, "")
	resp := req(t, app, "PATCH", "/api/v1/book/999", "", map[string]any{"title": "X", "rating": 1})
	if resp.StatusCode != 404 {
		t.Fatalf("want 404, got %d", resp.StatusCode)
	}
	list := decode[[]any](t, req(t, app, "GET", "/api/v1/book", "", nil))
	if len(list) != 0 {
		t.Fatalf("phantom record created: %d", len(list))
	}
}

func TestPagination(t *testing.T) {
	app := newTestApp(t, "")
	for _, title := range []string{"a", "b", "c"} {
		req(t, app, "POST", "/api/v1/book", "", map[string]any{"title": title, "rating": 1})
	}
	page := decode[[]any](t, req(t, app, "GET", "/api/v1/book?limit=2", "", nil))
	if len(page) != 2 {
		t.Fatalf("limit=2 want 2, got %d", len(page))
	}
	page2 := decode[[]any](t, req(t, app, "GET", "/api/v1/book?limit=2&offset=2", "", nil))
	if len(page2) != 1 {
		t.Fatalf("offset=2 want 1, got %d", len(page2))
	}
}

func TestAuth(t *testing.T) {
	app := newTestApp(t, "secret")

	// Reads are public.
	if r := req(t, app, "GET", "/api/v1/book", "", nil); r.StatusCode != 200 {
		t.Fatalf("public read want 200, got %d", r.StatusCode)
	}
	// Write without key -> 401.
	if r := req(t, app, "POST", "/api/v1/book", "", map[string]any{"title": "X", "rating": 1}); r.StatusCode != 401 {
		t.Fatalf("no key want 401, got %d", r.StatusCode)
	}
	// Write with wrong key -> 401.
	if r := req(t, app, "POST", "/api/v1/book", "nope", map[string]any{"title": "X", "rating": 1}); r.StatusCode != 401 {
		t.Fatalf("wrong key want 401, got %d", r.StatusCode)
	}
	// Write with correct key -> 201.
	if r := req(t, app, "POST", "/api/v1/book", "secret", map[string]any{"title": "X", "rating": 1}); r.StatusCode != 201 {
		t.Fatalf("valid key want 201, got %d", r.StatusCode)
	}
}

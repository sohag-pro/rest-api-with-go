package book

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"path/filepath"
	"testing"

	"restapi/database"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupApp(t *testing.T) *fiber.App {
	t.Helper()

	dbPath := filepath.Join(t.TempDir(), "test.db")
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := db.AutoMigrate(&Book{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	database.DBConn = db

	app := fiber.New()
	app.Get("/api/v1/book", GetBooks)
	app.Get("/api/v1/book/:id", GetBook)
	app.Post("/api/v1/book", NewBook)
	app.Patch("/api/v1/book/:id", UpdateBook)
	app.Delete("/api/v1/book/:id", DeleteBooks)
	return app
}

func doJSON(t *testing.T, app *fiber.App, method, url string, body any) *http.Response {
	t.Helper()
	var rdr io.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		rdr = bytes.NewReader(b)
	}
	req, _ := http.NewRequest(method, url, rdr)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("%s %s: %v", method, url, err)
	}
	return resp
}

func decodeBook(t *testing.T, resp *http.Response) Book {
	t.Helper()
	var b Book
	if err := json.NewDecoder(resp.Body).Decode(&b); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return b
}

func TestCreateBook(t *testing.T) {
	app := setupApp(t)
	resp := doJSON(t, app, "POST", "/api/v1/book", Book{Title: "Go", Author: "X", Rating: 5})
	if resp.StatusCode != 200 {
		t.Fatalf("want 200, got %d", resp.StatusCode)
	}
	b := decodeBook(t, resp)
	if b.ID == 0 || b.Title != "Go" {
		t.Fatalf("unexpected book: %+v", b)
	}
}

func TestCreateBookMissingTitle(t *testing.T) {
	app := setupApp(t)
	resp := doJSON(t, app, "POST", "/api/v1/book", Book{Author: "X", Rating: 3})
	if resp.StatusCode != 400 {
		t.Fatalf("want 400, got %d", resp.StatusCode)
	}
}

func TestCreateBookBadRating(t *testing.T) {
	app := setupApp(t)
	resp := doJSON(t, app, "POST", "/api/v1/book", Book{Title: "Go", Rating: 9})
	if resp.StatusCode != 400 {
		t.Fatalf("want 400, got %d", resp.StatusCode)
	}
}

func TestCreateBookTrimsWhitespaceTitle(t *testing.T) {
	app := setupApp(t)
	resp := doJSON(t, app, "POST", "/api/v1/book", Book{Title: "   ", Rating: 1})
	if resp.StatusCode != 400 {
		t.Fatalf("want 400 for whitespace title, got %d", resp.StatusCode)
	}
}

func TestGetBookNotFound(t *testing.T) {
	app := setupApp(t)
	resp := doJSON(t, app, "GET", "/api/v1/book/999", nil)
	if resp.StatusCode != 404 {
		t.Fatalf("want 404, got %d", resp.StatusCode)
	}
}

func TestUpdateBookNotFound(t *testing.T) {
	app := setupApp(t)
	resp := doJSON(t, app, "PATCH", "/api/v1/book/999", Book{Title: "Go", Rating: 1})
	if resp.StatusCode != 404 {
		t.Fatalf("want 404, got %d", resp.StatusCode)
	}
	// Must not have created a phantom record.
	list := doJSON(t, app, "GET", "/api/v1/book", nil)
	var books []Book
	json.NewDecoder(list.Body).Decode(&books)
	if len(books) != 0 {
		t.Fatalf("phantom record created: %d books", len(books))
	}
}

func TestDeleteBookNotFound(t *testing.T) {
	app := setupApp(t)
	resp := doJSON(t, app, "DELETE", "/api/v1/book/999", nil)
	if resp.StatusCode != 404 {
		t.Fatalf("want 404, got %d", resp.StatusCode)
	}
}

func TestCRUDFlow(t *testing.T) {
	app := setupApp(t)

	created := decodeBook(t, doJSON(t, app, "POST", "/api/v1/book", Book{Title: "A", Author: "B", Rating: 3}))
	id := "1"
	if created.ID != 1 {
		t.Fatalf("want ID 1, got %d", created.ID)
	}

	got := decodeBook(t, doJSON(t, app, "GET", "/api/v1/book/"+id, nil))
	if got.Title != "A" {
		t.Fatalf("want title A, got %q", got.Title)
	}

	updated := decodeBook(t, doJSON(t, app, "PATCH", "/api/v1/book/"+id, Book{Title: "A2", Author: "B2", Rating: 4}))
	if updated.Title != "A2" || updated.Rating != 4 {
		t.Fatalf("update failed: %+v", updated)
	}

	if del := doJSON(t, app, "DELETE", "/api/v1/book/"+id, nil); del.StatusCode != 200 {
		t.Fatalf("delete want 200, got %d", del.StatusCode)
	}

	if after := doJSON(t, app, "GET", "/api/v1/book/"+id, nil); after.StatusCode != 404 {
		t.Fatalf("after delete want 404, got %d", after.StatusCode)
	}
}

func TestPagination(t *testing.T) {
	app := setupApp(t)
	for _, title := range []string{"one", "two", "three"} {
		doJSON(t, app, "POST", "/api/v1/book", Book{Title: title, Rating: 1})
	}

	resp := doJSON(t, app, "GET", "/api/v1/book?limit=2&offset=0", nil)
	var page []Book
	json.NewDecoder(resp.Body).Decode(&page)
	if len(page) != 2 {
		t.Fatalf("limit=2 want 2 books, got %d", len(page))
	}

	resp2 := doJSON(t, app, "GET", "/api/v1/book?limit=2&offset=2", nil)
	var page2 []Book
	json.NewDecoder(resp2.Body).Decode(&page2)
	if len(page2) != 1 {
		t.Fatalf("offset=2 want 1 book, got %d", len(page2))
	}
}

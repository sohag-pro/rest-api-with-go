# Books REST API (Go + Fiber + GORM)

A small, clean REST API for managing a collection of books. Built with
[Fiber](https://gofiber.io/) (fast HTTP framework) and [GORM](https://gorm.io/)
(ORM) on top of SQLite. Useful as a starter/reference for structuring a Go REST
service with separated routing, handlers, and database packages.

## Features

- Full CRUD for books (create, read, list, update, delete)
- SQLite storage with automatic schema migration on startup
- Soft deletes and automatic timestamps via `gorm.Model`
- Package-separated layout (`main`, `book`, `database`)

## Tech Stack

| Component | Library | Version |
|-----------|---------|---------|
| HTTP framework | `github.com/gofiber/fiber/v2` | v2.52.13 |
| ORM | `gorm.io/gorm` | v1.31.1 |
| DB driver | `gorm.io/driver/sqlite` | v1.6.0 |
| Language | Go | 1.24+ |

## Project Layout

```
.
├── main.go             # entry point: Fiber app, DB init/migrate, server start
├── route.go            # route → handler wiring
├── database/
│   └── database.go     # shared *gorm.DB connection
├── book/
│   └── book.go         # Book model + CRUD handlers
├── go.mod / go.sum     # module + dependencies
└── books.db            # SQLite file (auto-created, gitignored)
```

## Prerequisites

- Go **1.24+** (`go version`)
- A C compiler — the SQLite driver (`mattn/go-sqlite3`) uses cgo, so `CGO_ENABLED=1`
  is required (default on macOS/Linux with Xcode CLT or gcc installed)

## Getting Started

```bash
# clone
git clone https://github.com/shuvopro/rest-api-with-go.git
cd rest-api-with-go

# download dependencies
go mod download

# run (server listens on :3000)
go run .
```

On first run the app creates `books.db` and migrates the `books` table.

### Build a binary

```bash
go build -o books-api .
./books-api
```

## Configuration

Currently fixed values (edit source to change):

| Setting | Value | Location |
|---------|-------|----------|
| Listen port | `:3000` | `main.go` |
| SQLite file | `books.db` | `main.go` |

## API Reference

Base path: `http://localhost:3000/api/v1`

| Method | Endpoint     | Description        | Body            |
|--------|--------------|--------------------|-----------------|
| GET    | `/book`      | List all books     | —               |
| GET    | `/book/:id`  | Get one book by ID | —               |
| POST   | `/book`      | Create a book      | Book JSON       |
| PATCH  | `/book/:id`  | Update a book      | Book JSON       |
| DELETE | `/book/:id`  | Delete a book      | —               |

### Book object

```json
{
  "title": "Clean Code",
  "author": "Robert C. Martin",
  "rating": 5
}
```

Responses also include GORM metadata: `ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt`.
`title` is required on create (returns `400` if missing).

### Status codes

| Code | Meaning |
|------|---------|
| 200  | Success |
| 400  | Missing required field (`title`) |
| 404  | Book not found |
| 406  | Request body could not be parsed |

## Example Requests

```bash
# Create
curl -X POST localhost:3000/api/v1/book \
  -H 'Content-Type: application/json' \
  -d '{"title":"Clean Code","author":"Robert C. Martin","rating":5}'

# List
curl localhost:3000/api/v1/book

# Get one
curl localhost:3000/api/v1/book/1

# Update
curl -X PATCH localhost:3000/api/v1/book/1 \
  -H 'Content-Type: application/json' \
  -d '{"title":"Clean Code","author":"Uncle Bob","rating":4}'

# Delete
curl -X DELETE localhost:3000/api/v1/book/1
```

## Development

```bash
go build ./...   # compile
go vet ./...     # static checks
gofmt -l .       # list unformatted files (empty = clean)
```

## Known Limitations

- No tests yet
- Config is hardcoded (port, DB path) — no env/flag support
- No pagination, auth, or input validation beyond `title`

## License

No license specified. Add one (e.g. MIT) before public distribution.

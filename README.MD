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
├── route.go            # route → handler wiring (groups + auth)
├── config.go           # env-var configuration loader
├── middleware.go       # API-key auth middleware
├── main_test.go        # config + auth middleware tests
├── database/
│   └── database.go     # shared *gorm.DB connection
├── book/
│   ├── book.go         # Book model + CRUD handlers
│   └── book_test.go    # handler/validation/pagination tests
├── .env.example        # environment variable template
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

# set up environment (optional — defaults work without it)
cp .env.example .env
export $(grep -v '^#' .env | xargs)

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

Configured via environment variables (all optional):

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT`     | `3000`     | Server listen port |
| `DB_PATH`  | `books.db` | SQLite file path |
| `API_KEY`  | *(empty)*  | If set, write endpoints require this key. Empty disables auth. |

```bash
PORT=8080 DB_PATH=/data/books.db API_KEY=s3cret go run .
```

See [`.env.example`](.env.example) for a template. The app reads real environment
variables (it does not auto-load `.env`), so export them first:

```bash
cp .env.example .env
export $(grep -v '^#' .env | xargs) && go run .
```

## Authentication

Read endpoints (`GET`) are always public. Write endpoints (`POST`, `PATCH`,
`DELETE`) are guarded only when `API_KEY` is set. Authenticated requests must send
the key in the `X-API-Key` header:

```bash
curl -X POST localhost:3000/api/v1/book \
  -H 'X-API-Key: s3cret' \
  -H 'Content-Type: application/json' \
  -d '{"title":"Go","rating":5}'
```

Without `API_KEY` set, auth is disabled and all endpoints are open.

## API Reference

Base path: `http://localhost:3000/api/v1`

| Method | Endpoint     | Description        | Auth | Body            |
|--------|--------------|--------------------|------|-----------------|
| GET    | `/book`      | List books (paginated) | — | —           |
| GET    | `/book/:id`  | Get one book by ID | —    | —               |
| POST   | `/book`      | Create a book      | key  | Book JSON       |
| PATCH  | `/book/:id`  | Update a book      | key  | Book JSON       |
| DELETE | `/book/:id`  | Delete a book      | key  | —               |

### Pagination

`GET /book` accepts query parameters:

| Param    | Default | Notes |
|----------|---------|-------|
| `limit`  | `10`    | Max `100`; out-of-range values fall back to `10` |
| `offset` | `0`     | Negative values clamp to `0` |

```bash
curl 'localhost:3000/api/v1/book?limit=20&offset=40'
```

### Book object

```json
{
  "title": "Clean Code",
  "author": "Robert C. Martin",
  "rating": 5
}
```

Responses also include GORM metadata: `ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt`.

### Validation

- `title` — required, whitespace-trimmed, must be non-empty
- `rating` — must be between `0` and `5`

### Status codes

| Code | Meaning |
|------|---------|
| 200  | Success |
| 400  | Validation failed (missing `title` or `rating` out of range) |
| 401  | Missing/invalid `X-API-Key` (when auth enabled) |
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
go test ./...    # run tests
gofmt -l .       # list unformatted files (empty = clean)
```

## Testing

Tests cover handler CRUD, validation, pagination, not-found paths, and the
API-key middleware. Each test uses an isolated temporary SQLite file, so no
running server or shared state is required.

```bash
go test ./... -v
```

## License

No license specified. Add one (e.g. MIT) before public distribution.

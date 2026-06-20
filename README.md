# Books REST API

[![CI](https://github.com/sohag-pro/rest-api-with-go/actions/workflows/ci.yml/badge.svg)](https://github.com/sohag-pro/rest-api-with-go/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/badge/go-1.24-00ADD8?logo=go)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

A production-grade REST API for managing a collection of books, built in Go with
[Fiber](https://gofiber.io/) and [GORM](https://gorm.io/) on SQLite. It
demonstrates a clean, layered architecture, graceful shutdown, structured
logging, health checks, containerization, and CI.

## Features

- **Layered architecture** — `handler → service → repository` with interfaces for testability
- **Graceful shutdown** — listens for `SIGINT`/`SIGTERM`, drains in-flight requests with a bounded timeout
- **Structured logging** — JSON logs via `log/slog`, per-request log lines with request IDs
- **Health checks** — `/healthz` probe that pings the database
- **Middleware** — request ID, panic recovery, CORS, request logging
- **Optional API-key auth** — guards write endpoints when configured
- **Input validation & pagination** — on all relevant endpoints
- **Consistent JSON error envelope** — uniform `{error, code}` responses
- **Config from environment** — validated on startup
- **Containerized** — multi-stage Dockerfile + docker-compose
- **CI** — GitHub Actions: format check, vet, race tests, lint, Docker build
- **OpenAPI spec** — see [`api/openapi.yaml`](api/openapi.yaml)

## Tech Stack

| Component | Library | Version |
|-----------|---------|---------|
| HTTP framework | `github.com/gofiber/fiber/v2` | v2.52.13 |
| ORM | `gorm.io/gorm` | v1.31.1 |
| DB driver | `gorm.io/driver/sqlite` | v1.6.0 |
| Logging | `log/slog` (stdlib) | — |
| Language | Go | 1.24+ |

## Architecture

Requests flow through clearly separated layers, each depending only on the one
below it via interfaces:

```
HTTP request
   │
   ▼
┌──────────┐   parse, map errors, status codes
│ Handler  │   internal/book/handler.go
└────┬─────┘
     ▼
┌──────────┐   business rules: normalize, validate, orchestrate
│ Service  │   internal/book/service.go
└────┬─────┘
     ▼
┌──────────┐   persistence (Repository interface)
│Repository│   internal/book/repository.go  →  GORM / SQLite
└──────────┘
```

The `Repository` interface lets the service be unit-tested with an in-memory fake
(no database required), while the full stack is covered by integration tests.

## Project Layout

```
.
├── cmd/
│   └── server/
│       └── main.go         # entry point: config, logger, DB, server, graceful shutdown
├── internal/
│   ├── config/             # env-var loading + validation
│   ├── database/           # GORM/SQLite open, pool tuning, ping, migrate
│   ├── book/               # domain: model, repository, service, handler
│   ├── server/             # Fiber app builder + middleware
│   └── response/           # JSON error envelope
├── api/
│   └── openapi.yaml        # OpenAPI 3 spec
├── .github/workflows/ci.yml
├── Dockerfile
├── docker-compose.yml
├── Makefile
├── .golangci.yml
├── .env.example
└── go.mod / go.sum
```

## Prerequisites

- Go **1.24+** (`go version`)
- A C compiler — the SQLite driver (`mattn/go-sqlite3`) uses cgo, so
  `CGO_ENABLED=1` (default on macOS/Linux with Xcode CLT or gcc)
- Optional: Docker, `make`, `golangci-lint`

## Quick Start

### Local

```bash
git clone https://github.com/sohag-pro/rest-api-with-go.git
cd rest-api-with-go

# optional: set up environment (defaults work without it)
cp .env.example .env
export $(grep -v '^#' .env | xargs)

make run            # or: go run ./cmd/server
```

The server listens on `:3000`, creates the SQLite file, and migrates the schema
on first run.

### Docker

```bash
docker compose up --build
# or
make docker && docker run -p 3000:3000 books-api:latest
```

## Configuration

All settings come from environment variables (validated on startup). See
[`.env.example`](.env.example).

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT`             | `3000`     | Listen port (1–65535) |
| `DB_PATH`          | `books.db` | SQLite file path |
| `API_KEY`          | *(empty)*  | If set, write endpoints require it; empty disables auth |
| `LOG_LEVEL`        | `info`     | `debug` \| `info` \| `warn` \| `error` |
| `READ_TIMEOUT`     | `10s`      | HTTP read timeout |
| `WRITE_TIMEOUT`    | `10s`      | HTTP write timeout |
| `SHUTDOWN_TIMEOUT` | `10s`      | Graceful shutdown grace period |

## Authentication

Read endpoints (`GET`) are always public. Write endpoints (`POST`, `PATCH`,
`DELETE`) require the `X-API-Key` header **only when `API_KEY` is set**:

```bash
curl -X POST localhost:3000/api/v1/book \
  -H 'X-API-Key: s3cret' \
  -H 'Content-Type: application/json' \
  -d '{"title":"Go","rating":5}'
```

## API Reference

Base path: `http://localhost:3000`

| Method | Endpoint            | Description             | Auth |
|--------|---------------------|-------------------------|------|
| GET    | `/healthz`          | Liveness/readiness probe | —   |
| GET    | `/api/v1/book`      | List books (paginated)  | —    |
| GET    | `/api/v1/book/:id`  | Get one book by ID      | —    |
| POST   | `/api/v1/book`      | Create a book           | key  |
| PATCH  | `/api/v1/book/:id`  | Update a book           | key  |
| DELETE | `/api/v1/book/:id`  | Delete a book           | key  |

Full schema in [`api/openapi.yaml`](api/openapi.yaml).

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

- `title` — required, whitespace-trimmed, non-empty
- `rating` — integer between `0` and `5`

### Pagination

`GET /api/v1/book` accepts query parameters:

| Param    | Default | Notes |
|----------|---------|-------|
| `limit`  | `10`    | Max `100`; out-of-range falls back to `10` |
| `offset` | `0`     | Negative clamps to `0` |

```bash
curl 'localhost:3000/api/v1/book?limit=20&offset=40'
```

### Error envelope

All errors return a consistent JSON shape:

```json
{ "error": "book not found", "code": 404 }
```

| Code | Meaning |
|------|---------|
| 200  | Success |
| 201  | Created |
| 400  | Validation failed / invalid id |
| 401  | Missing/invalid `X-API-Key` (when auth enabled) |
| 404  | Not found |
| 406  | Request body could not be parsed |
| 500  | Internal error |
| 503  | Database unavailable (`/healthz`) |

## Example Requests

```bash
# Create
curl -X POST localhost:3000/api/v1/book \
  -H 'Content-Type: application/json' \
  -d '{"title":"Clean Code","author":"Robert C. Martin","rating":5}'

# List (paginated)
curl 'localhost:3000/api/v1/book?limit=10&offset=0'

# Get one
curl localhost:3000/api/v1/book/1

# Update
curl -X PATCH localhost:3000/api/v1/book/1 \
  -H 'Content-Type: application/json' \
  -d '{"title":"Clean Code","author":"Uncle Bob","rating":4}'

# Delete
curl -X DELETE localhost:3000/api/v1/book/1

# Health
curl localhost:3000/healthz
```

## Development

```bash
make help     # list all targets
make run      # run the server
make test     # race-enabled tests
make cover    # coverage report (HTML)
make vet      # go vet
make fmt      # gofmt -w
make lint     # golangci-lint
make build    # build ./bin/books-api
make docker   # build Docker image
```

## Testing

Tests cover service logic (with an in-memory fake repository), full HTTP
integration (handlers → service → repository), validation, pagination,
auth, health checks, and config loading. Each test uses an isolated temporary
SQLite file — no running server or shared state required.

```bash
go test -race ./...
```

## License

[MIT](LICENSE) © Sohag Hasan

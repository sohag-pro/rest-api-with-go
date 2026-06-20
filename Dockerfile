# syntax=docker/dockerfile:1

# --- Build stage ---
FROM golang:1.24-alpine AS builder

# cgo is required by the SQLite driver (mattn/go-sqlite3).
RUN apk add --no-cache gcc musl-dev

WORKDIR /src

# Cache dependencies first.
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/books-api ./cmd/server

# --- Runtime stage ---
FROM alpine:3.20

RUN apk add --no-cache ca-certificates wget \
    && adduser -D -u 10001 app \
    && mkdir -p /data && chown app /data

WORKDIR /app
COPY --from=builder /out/books-api .

USER app
EXPOSE 3000
ENV PORT=3000 \
    DB_PATH=/data/books.db
VOLUME ["/data"]

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget -qO- "http://localhost:${PORT}/healthz" || exit 1

ENTRYPOINT ["./books-api"]

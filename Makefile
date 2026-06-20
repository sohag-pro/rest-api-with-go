.DEFAULT_GOAL := help
BINARY := books-api
PKG := ./...

.PHONY: help
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
	  awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-12s\033[0m %s\n", $$1, $$2}'

.PHONY: run
run: ## Run the server
	go run ./cmd/server

.PHONY: build
build: ## Build the binary into ./bin
	CGO_ENABLED=1 go build -trimpath -ldflags="-s -w" -o bin/$(BINARY) ./cmd/server

.PHONY: test
test: ## Run tests with race detector
	go test -race $(PKG)

.PHONY: cover
cover: ## Run tests and open coverage report
	go test -coverprofile=coverage.out $(PKG)
	go tool cover -html=coverage.out

.PHONY: vet
vet: ## Run go vet
	go vet $(PKG)

.PHONY: fmt
fmt: ## Format code
	gofmt -w .

.PHONY: lint
lint: ## Run golangci-lint (must be installed)
	golangci-lint run

.PHONY: tidy
tidy: ## Tidy module dependencies
	go mod tidy

.PHONY: docker
docker: ## Build the Docker image
	docker build -t $(BINARY):latest .

.PHONY: clean
clean: ## Remove build artifacts
	rm -rf bin coverage.out

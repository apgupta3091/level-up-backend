-include .env
export

BINARY_NAME    := levelup
CMD_PATH       := ./cmd/server
MIGRATE_PATH   := ./db/migrations
DB_URL         := postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable

.PHONY: help
help: ## Display this help message
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

# ── Development ──────────────────────────────────────────────────────────────

.PHONY: run
run: ## Run the server
	go run $(CMD_PATH)/main.go

.PHONY: build
build: ## Build binary to ./bin/levelup
	mkdir -p bin
	go build -o bin/$(BINARY_NAME) $(CMD_PATH)/main.go

.PHONY: tidy
tidy: ## Tidy go.mod and go.sum
	go mod tidy

# ── Database ─────────────────────────────────────────────────────────────────

.PHONY: db/up
db/up: ## Start postgres via docker compose
	docker compose up postgres -d

.PHONY: db/down
db/down: ## Stop postgres container
	docker compose down postgres

.PHONY: db/shell
db/shell: ## Open psql shell
	docker compose exec postgres psql -U $(POSTGRES_USER) -d $(POSTGRES_DB)

# ── Migrations ────────────────────────────────────────────────────────────────

.PHONY: migrate/up
migrate/up: ## Apply all pending migrations
	migrate -database "$(DB_URL)" -path $(MIGRATE_PATH) up

.PHONY: migrate/down
migrate/down: ## Roll back last migration
	migrate -database "$(DB_URL)" -path $(MIGRATE_PATH) down 1

.PHONY: migrate/down/all
migrate/down/all: ## Roll back ALL migrations (destructive!)
	migrate -database "$(DB_URL)" -path $(MIGRATE_PATH) down -all

.PHONY: migrate/status
migrate/status: ## Show current migration version
	migrate -database "$(DB_URL)" -path $(MIGRATE_PATH) version

.PHONY: migrate/create
migrate/create: ## Create migration pair. Usage: make migrate/create name=add_refresh_tokens
	@if [ -z "$(name)" ]; then echo "Usage: make migrate/create name=<name>"; exit 1; fi
	migrate create -ext sql -dir $(MIGRATE_PATH) -seq $(name)

.PHONY: migrate/force
migrate/force: ## Force migration version. Usage: make migrate/force version=7
	migrate -database "$(DB_URL)" -path $(MIGRATE_PATH) force $(version)

# ── SQLC ─────────────────────────────────────────────────────────────────────

.PHONY: sqlc/generate
sqlc/generate: ## Regenerate SQLC database code
	sqlc generate

.PHONY: sqlc/verify
sqlc/verify: ## Verify SQLC queries without generating
	sqlc vet

# ── Stripe ───────────────────────────────────────────────────────────────────

.PHONY: stripe/listen
stripe/listen: ## Forward Stripe webhook events to local server
	stripe listen --forward-to localhost:$(PORT)/payments/webhook

# ── Quality ──────────────────────────────────────────────────────────────────

.PHONY: test
test: ## Run all tests with race detector
	go test ./... -v -race -timeout 30s

.PHONY: test/short
test/short: ## Run tests, skip integration tests
	go test ./... -v -short

.PHONY: lint
lint: ## Run golangci-lint
	golangci-lint run ./...

.PHONY: fmt
fmt: ## Format all Go code
	gofmt -w .

.PHONY: vet
vet: ## Run go vet
	go vet ./...

# ── Docker ───────────────────────────────────────────────────────────────────

.PHONY: docker/up
docker/up: ## Start all services (postgres + mailhog)
	docker compose up -d

.PHONY: docker/down
docker/down: ## Stop all services
	docker compose down

.PHONY: docker/logs
docker/logs: ## Tail logs from all compose services
	docker compose logs -f

# ── Setup ────────────────────────────────────────────────────────────────────

.PHONY: setup
setup: ## First-time dev setup
	cp -n .env.example .env || true
	docker compose up postgres mailhog -d
	sleep 3
	$(MAKE) migrate/up
	$(MAKE) sqlc/generate
	@echo "Setup complete. Run 'make run' to start the server."

# Makefile for Contact Management Go Project

# Variables
APP_NAME=contact-management
BUILD_DIR=./bin
BINARY=$(BUILD_DIR)/$(APP_NAME)
MAIN_FILE=./main.go
MIGRATIONS_DIR=./src/migrations

# Database configuration (load from .env or set defaults)
include .env
export

# Migration targets
.PHONY: migrate migrate-rollback migrate-fresh migrate-up migrate-down

migrate:
	@echo "Running all migrations up..."
	migrate -path $(MIGRATIONS_DIR) -database "mysql://$(DB_USER):$(DB_PASSWORD)@tcp($(DB_HOST):$(DB_PORT))/$(DB_NAME)" up

migrate-rollback:
	@echo "Rolling back last migration..."
	migrate -path $(MIGRATIONS_DIR) -database "mysql://$(DB_USER):$(DB_PASSWORD)@tcp($(DB_HOST):$(DB_PORT))/$(DB_NAME)" down 1

migrate-down:
	@echo "Rolling back all migrations..."
	migrate -path $(MIGRATIONS_DIR) -database "mysql://$(DB_USER):$(DB_PASSWORD)@tcp($(DB_HOST):$(DB_PORT))/$(DB_NAME)" down

migrate-fresh:
	@echo "Refreshing migrations (down then up)..."
	@$(MAKE) migrate-down
	@$(MAKE) migrate

# Build targets
.PHONY: build clean

build: clean
	@echo "Building application..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BINARY) $(MAIN_FILE)
	@echo "Build complete: $(BINARY)"

clean:
	@echo "Cleaning old builds..."
	@rm -rf $(BUILD_DIR)
	@echo "Clean complete"

# Run targets
.PHONY: run run-dev

run: build
	@echo "Running application from binary..."
	$(BINARY)

run-dev:
	@echo "Running application from source (main.go)..."
	go run $(MAIN_FILE)

# Development helpers
.PHONY: test deps tidy

test:
	@echo "Running tests..."
	go test ./test/... -v

deps:
	@echo "Installing dependencies..."
	go mod download

tidy:
	@echo "Tidying go.mod..."
	go mod tidy

# Help
.PHONY: help

help:
	@echo "Available targets:"
	@echo "  migrate           - Run all migrations up"
	@echo "  migrate-rollback  - Rollback last migration (down 1)"
	@echo "  migrate-down      - Rollback all migrations"
	@echo "  migrate-fresh     - Down all then up all migrations"
	@echo "  build             - Build the application (removes old build first)"
	@echo "  clean             - Remove build directory"
	@echo "  run               - Build and run from binary"
	@echo "  run-dev           - Run directly from main.go"
	@echo "  test              - Run all tests"
	@echo "  deps              - Download dependencies"
	@echo "  tidy              - Tidy go.mod"
	@echo "  help              - Show this help message"

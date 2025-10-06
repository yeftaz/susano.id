.PHONY: help dev build migrate migrate-up migrate-down migrate-fresh seed test test-unit test-integration routes clean

# Load environment variables from .env file
ifneq (,$(wildcard .env))
    include .env
    export
endif

# Variables (with fallback defaults if .env not found)
APP_NAME ?= susano-api
DB_URL ?= postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)
MIGRATION_PATH = internal/database/migrations

help:
	@echo "Available commands:"
	@echo "  make dev              - Run development server with hot reload"
	@echo "  make build            - Build production binary"
	@echo "  make migrate          - Run all pending migrations"
	@echo "  make migrate-up       - Run next migration"
	@echo "  make migrate-down     - Rollback last migration"
	@echo "  make migrate-fresh    - Drop all tables and re-migrate"
	@echo "  make seed             - Run database seeders"
	@echo "  make routes           - Show all API routes"
	@echo "  make test             - Run all tests"
	@echo "  make test-unit        - Run unit tests only"
	@echo "  make test-integration - Run integration tests only"
	@echo "  make clean            - Clean build artifacts and logs"

dev:
	@echo "Starting development server..."
	@go run cmd/api/main.go

build:
	@echo "Building production binary..."
	@go build -o bin/$(APP_NAME) cmd/api/main.go
	@echo "Binary created at bin/$(APP_NAME)"

migrate:
	@echo "Running migrations..."
	@migrate -path $(MIGRATION_PATH) -database "$(DB_URL)" up

migrate-up:
	@echo "Running next migration..."
	@migrate -path $(MIGRATION_PATH) -database "$(DB_URL)" up 1

migrate-down:
	@echo "Rolling back last migration..."
	@migrate -path $(MIGRATION_PATH) -database "$(DB_URL)" down 1

migrate-fresh:
	@echo "Dropping all tables and re-migrating..."
	@migrate -path $(MIGRATION_PATH) -database "$(DB_URL)" drop -f
	@$(MAKE) migrate
	@echo "Fresh migration completed!"

seed:
	@echo "Running seeders..."
	@psql "$(DB_URL)" -f internal/database/seeds/database_seeder.sql
	@echo "Seeding completed!"

routes:
	@echo "Generating API routes..."
	@go run cmd/api/main.go -routes

test:
	@echo "Running all tests..."
	@go test -v -race -cover ./...

test-unit:
	@echo "Running unit tests..."
	@go test -v -race -cover ./tests/unit/...

test-integration:
	@echo "Running integration tests..."
	@go test -v -race -cover ./tests/integration/...

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@rm -f storage/logs/*.log
	@echo "Clean completed!"

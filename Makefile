SHELL := /bin/bash

# Configurable variables
IMAGE_NAME ?= todo-app
IMAGE_TAG  ?= latest
BUILD_TARGET ?= builder
COMPOSE_DEV  ?= docker-compose.dev.yml
COMPOSE_PROD ?= docker-compose.prod.yml

.PHONY: help build build-prod extract-bin up-dev up-prod down logs clean fmt vet test sqlc-gen schema seed docker-push

help:
	@echo "Available targets:"
	@echo "  build         Build the builder image ($(IMAGE_NAME):$(IMAGE_TAG))"
	@echo "  build-prod    Build the production image and extract binary"
	@echo "  extract-bin   Extract compiled binary from builder image"
	@echo "  up-dev        Start development stack (docker compose dev)"
	@echo "  up-prod       Start production stack (docker compose prod, detached)"
	@echo "  down          Stop stacks (dev + prod)"
	@echo "  logs          Follow dev compose logs"
	@echo "  fmt           Run 'go fmt'"
	@echo "  vet           Run 'go vet'"
	@echo "  test          Run 'go test'"
	@echo "  sqlc-gen      Run 'sqlc generate'"
	@echo "  schema        Apply SQL schema to DB in docker"
	@echo "  seed          Run seed script (go run ./cmd/seed)"
	@echo "  docker-push   Push image to registry"

build:
	@echo "Building image ($(IMAGE_NAME):$(IMAGE_TAG))..."
	docker build --target $(BUILD_TARGET) -t $(IMAGE_NAME):$(IMAGE_TAG) .

extract-bin:
	@echo "Extracting binary from builder image..."
	docker create --name tmp $(IMAGE_NAME):$(IMAGE_TAG) || true
	docker cp tmp:/app/main ./main || true
	docker rm -f tmp || true

build-prod: build extract-bin
	@echo "Building production image from extracted binary..."
	docker build --target prod -t $(IMAGE_NAME):prod .

up-dev:
	docker compose -f $(COMPOSE_DEV) up --build

up-prod:
	docker compose -f $(COMPOSE_PROD) up --build -d

down:
	docker compose -f $(COMPOSE_DEV) down || true
	docker compose -f $(COMPOSE_PROD) down || true

logs:
	docker compose -f $(COMPOSE_DEV) logs -f

clean:
	@echo "Cleaning local artifacts..."
	rm -f ./main || true

fmt:
	go fmt ./...

vet:
	go vet ./...

test:
	go test ./...

# ----------------------------
# Database / sqlc targets
# ----------------------------
DB_USER ?= user
DB_NAME ?= todo_db
DB_HOST ?= db

sqlc-gen:
	sqlc generate

schema:
	docker compose -f $(COMPOSE_DEV) exec -T $(DB_HOST) psql -U $(DB_USER) -d $(DB_NAME) < sql/schema/schema.sql

seed:
	go run ./cmd/seed

docker-push:
	docker push $(IMAGE_NAME):$(IMAGE_TAG)
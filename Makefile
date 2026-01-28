SHELL := /bin/bash

.PHONY: help build build-prod up dev up-prod down logs

help:
	@echo "Available targets: build build-prod up-dev up-prod down logs"

build:
	@echo "Building production image..."
	docker build --target builder -t todo-app:builder .

build-prod: build
	@echo "Creating prod image from builder..."
	docker create --name tmp todo-app:builder || true
	docker cp tmp:/app/main ./main || true
	docker rm -f tmp || true
	docker build --target prod -t todo-app:prod .

up-dev:
	docker compose -f docker-compose.dev.yml up --build

up-prod:
	docker compose -f docker-compose.prod.yml up --build -d

down:
	docker compose -f docker-compose.dev.yml down || true
	docker compose -f docker-compose.prod.yml down || true

logs:
	docker compose -f docker-compose.dev.yml logs -f

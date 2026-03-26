.PHONY: run build tidy migrate-up migrate-down migrate-force migrate-version docker-up docker-down

# Load env variables
include .env
export

DB_URL=postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

# Server
run:
	go run cmd/server/main.go

build:
	go build -o bin/ledger-core cmd/server/main.go

tidy:
	go mod tidy

# Docker
docker-up:
	docker compose up -d

docker-down:
	docker compose down

# Migrations
migrate-up:
	migrate -path migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path migrations -database "$(DB_URL)" down 1

migrate-force:
	migrate -path migrations -database "$(DB_URL)" force $(version)

migrate-version:
	migrate -path migrations -database "$(DB_URL)" version

migrate-down-all:
	migrate -path migrations -database "$(DB_URL)" down
.PHONY: run build test tidy sqlc migrate-up migrate-down docker-up docker-down

# -----------------------------------------------------------------------
# Development
# -----------------------------------------------------------------------
run:
	go run ./cmd/server

build:
	go build -o bin/server ./cmd/server

test:
	go test -v -race ./...

tidy:
	go mod tidy

# -----------------------------------------------------------------------
# SQLC
# -----------------------------------------------------------------------
sqlc:
	sqlc generate

# -----------------------------------------------------------------------
# Migrations (set DB_URL in .env or export it before running)
# -----------------------------------------------------------------------
migrate-up:
	migrate -path db/migrations -database "mysql://$(DB_USER):$(DB_PASSWORD)@tcp($(DB_HOST):$(DB_PORT))/$(DB_NAME)" up

migrate-down:
	migrate -path db/migrations -database "mysql://$(DB_USER):$(DB_PASSWORD)@tcp($(DB_HOST):$(DB_PORT))/$(DB_NAME)" down 1

# -----------------------------------------------------------------------
# Docker
# -----------------------------------------------------------------------
docker-up:
	docker compose up --build

docker-down:
	docker compose down

docker-clean:
	docker compose down -v

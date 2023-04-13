lint:
	gofumpt -w .
	go mod tidy
	golangci-lint run

up:
	docker-compose up -d

down:
	docker-compose down

build:
	docker compose build --no-cache

.PHONY: build
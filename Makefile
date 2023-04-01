lint:
	gofumpt -w .
	go mod tidy
	golangci-lint run

run:
	docker-compose up -d

down:
	docker-compose down

build:
	docker compose build --no-cache

test:
	go test -v ./tests/main_test.go

.PHONY: build
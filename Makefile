APP_NAME=stin-smart-care-be

.PHONY: dev test lint migrate-up migrate-down seed

dev:
	go run ./cmd/api

test:
	go test ./...

lint:
	golangci-lint run

migrate-up:
	go run ./cmd/migrate -action up

migrate-down:
	go run ./cmd/migrate -action down

seed:
	go run ./cmd/seed

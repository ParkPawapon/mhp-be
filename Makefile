APP_NAME=stin-smart-care-be

.PHONY: dev test lint migrate-up migrate-down seed gen-jwt-secret test-integration

dev:
	go run ./cmd/api

test:
	go test ./...

test-integration:
	go test -tags=integration ./internal/repositories

lint:
	golangci-lint run

migrate-up:
	go run ./cmd/migrate -action up

migrate-down:
	go run ./cmd/migrate -action down

seed:
	go run ./cmd/seed

gen-jwt-secret:
	@python3 -c "import base64, secrets; print(base64.urlsafe_b64encode(secrets.token_bytes(64)).decode().rstrip('='))"

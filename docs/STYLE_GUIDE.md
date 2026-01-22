# STYLE GUIDE

## Go Style & Patterns
- Layering: handlers -> services -> repositories.
- Handlers: bind + validate + call service + map response.
- Services: business logic + transactions only.
- Repositories: database access only.
- DTOs: API request/response in `internal/models/dto`.
- DB models: `internal/models/db`.

## Context Usage
- `context.Context` is required everywhere.
- Handlers pass `c.Request.Context()` to services.
- Repositories accept `ctx context.Context` and use `WithContext` in GORM.

## Error Handling
- Wrap errors with `%w`.
- Use `domain.AppError` for stable codes and messages.
- Repositories map GORM errors to domain errors; never return raw GORM errors.
- Handlers call `httpx.Fail` for centralized mapping.

## DTO Mapping Rules
- Never return DB models directly to API.
- Map DB -> DTO in services or dedicated mappers.
- Apply data masking based on role (e.g., `citizen_id`).

## Validation Standard
- Use `go-playground/validator` for request validation.
- Return deterministic error details: map field -> message.
- Validation errors must return `VALIDATION_*` codes.

## Logging Standard
- Handler/middleware logs are structured with Zap.
- Required fields: `request_id`, `actor_id`, `role`, `route`, `method`, `status`, `latency_ms`, `ip`.
- Do not log secrets or PII (passwords, OTP, tokens).

## Testing Standard
- Unit tests named `*_test.go` with table-driven cases.
- Service tests use fakes/mocks for repositories.
- Handler tests use `httptest` and assert envelope.

## Transactions
- Transactions only in services using repository interfaces that accept `*gorm.DB` when needed.

## Background Jobs
- Jobs live in `internal/jobs` and call services (no direct handler logic).
- Use `NotificationService` for scheduling/sending; job interval is configured via env.
- Jobs must be idempotent and concurrency-safe (e.g., row locking with `FOR UPDATE SKIP LOCKED`).

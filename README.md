# STIN Smart Care Backend

Production backend for STIN Smart Care (Mobile + Web). This service uses Go, Gin, GORM, PostgreSQL, Redis, JWT, Prometheus, and OpenTelemetry scaffolding.

## Conventions
- `SYSTEM_CONVENTIONS.md` is the single source of truth for naming, API envelope, errors, and security.
- `docs/STYLE_GUIDE.md` covers code patterns and testing.
- `docs/API_CONTRACT.md` documents endpoints and envelopes.
- `docs/PRODUCTION_CHECKLIST.md` provides deployment and ops checklist.
- Module path: `github.com/ParkPawapon/mhp-be`.

## Requirements
- Go 1.22+
- PostgreSQL 14+
- Redis 7+
- `golangci-lint` (for `make lint`)

## Local Setup
1) Copy environment template:
   ```bash
   cp .env.example .env
   ```
2) Start dependencies:
   ```bash
   docker compose up -d postgres redis
   ```
3) Run migrations:
   ```bash
   make migrate-up
   ```
4) Run the API:
   ```bash
   make dev
   ```

## Migrations
- SQL migrations are the source of truth.
- Run `make migrate-up` and `make migrate-down`.

## Health & Metrics
- `GET /healthz`
- `GET /readyz`
- `GET /metrics`

## Base API
- `/api/v1`

## Seed
```bash
make seed
```

# SYSTEM CONVENTIONS (Single Source of Truth)

This document is the single source of truth for project conventions. All code must comply.

## Naming Standards
- Packages: lowercase, singular, no underscores (e.g., `handler`, `service`, `repository`).
- Files: lowercase, short and descriptive (e.g., `auth.go`, `user.go`).
- Types/Structs/Interfaces: PascalCase. Interfaces must end with `Repository` or `Service`.
- Methods: VerbNoun (e.g., `CreateAppointment`, `ListPatients`).
- Routes: kebab-case segments, plural resources, base `/api/v1`.
  - Good: `/patient-medicines`, `/auth/request-otp`, `/caregivers/assignments`
  - Bad: `/patientMedicines`, `/auth/requestOTP`
- DB: snake_case, plural table names, `*_id` for foreign keys.
- Enums: custom Go types with constants matching DB enums (e.g., `PATIENT`, `NURSE`).

## Folder Boundaries + Import Rules
- `handlers`: HTTP layer only. Bind/validate, call service, map response. NO DB calls.
- `services`: business logic + transactions. NO Gin context usage.
- `repositories`: persistence only. NO business rules.
- `models`: DB models and API DTOs are separate packages.
- `middleware`: request-id, logging, auth, rbac, recover. No business logic.
- `router`: route registrations only.
- `utils`: pure helpers only (no cross-layer imports).
- `config`, `database`, `cache`: wiring/infrastructure only.
- No circular imports.

## API Envelope Standard
All responses must follow the envelope:
- Success:
  ```json
  {"data":{},"meta":{"request_id":"..."}}
  ```
- Error:
  ```json
  {"error":{"code":"...","message":"...","details":null},"meta":{"request_id":"..."}}
  ```
- Pagination:
  ```json
  {"data":[],"meta":{"request_id":"...","page":1,"page_size":20,"total":123}}
  ```

Time format: RFC3339 in JSON. Store timestamps as UTC in DB (TIMESTAMPTZ).

## Error Code Taxonomy + HTTP Mapping
Codes are stable and mapped to HTTP:
- `AUTH_*` -> 401/403
- `USER_*` -> 404/409
- `MED_*` -> 400/404
- `APPT_*` -> 400/404
- `HEALTH_*` -> 400/404
- `CONTENT_*` -> 400/404
- `AUDIT_*` -> 400/404
- `RATE_*` -> 429
- `VALIDATION_*` -> 400
- `INTERNAL_*` -> 500

## Logging (Zap)
Every request log must include:
- `request_id`, `actor_id` (if exists), `role`, `route`, `method`, `status`, `latency_ms`, `ip`
Pattern: log once at end of request, structured JSON.

## Security Baseline
- Passwords: bcrypt.
- JWT: access 15m, refresh 30d; refresh rotation; revoke via Redis.
- OTP: 6 digits, TTL 5m; rate-limit per phone + IP via Redis.
- HTTPS required in production, HTTP allowed in local.
- CORS configurable via env for web admin origins.

## RBAC + Data Masking (PDPA-minded)
- PATIENT: self-only resources; no admin endpoints.
- CAREGIVER: read-only assigned patient data; never see `citizen_id`.
- NURSE: view patients; create appointments + notes.
- ADMIN: full access; publish content; audit logs.

Sensitive data rules:
- Never expose `password_hash`.
- `citizen_id` masked for non-admin.
- Only return minimum required fields for caregiver.

## Migrations
- SQL migrations are source of truth. No AutoMigrate in production path.
- Enable `pgcrypto`; create enums before tables; create indexes explicitly.
- Required indexes:
  - `users(username)`
  - `user_profiles(hn, citizen_id)`
  - `intake_history(user_id, target_date)`
  - `appointments(user_id, appt_datetime)`
  - `audit_logs(timestamp, actor_id)`

### updated_at Strategy
- Use GORM `autoUpdateTime` with UTC `NowFunc` (application-layer updates).

### Additional Tables & Enums
- Tables: `medicine_categories`, `medicine_category_items`, `device_tokens`, `notification_templates`, `notification_events`, `user_preferences`, `support_chat_requests`.
- Enum: `notification_status` = `PENDING`, `SENT`, `CANCELLED`, `FAILED`.
- `meal_timing` is a controlled string; valid options are documented in `docs/API_CONTRACT.md`.

## Versioning Rules
- Base API path `/api/v1`.
- Backwards compatibility: no breaking changes in v1 without new v2 path.

## Commit Message Convention
Use Conventional Commits:
- `feat: add otp rate limiting`
- `fix: handle expired otp`
- `chore: update deps`

## Do / Don't Examples
Do:
```go
if err := svc.CreateAppointment(ctx, req); err != nil {
  httpx.Fail(c, err)
  return
}
httpx.Created(c, resp)
```

Don't:
```go
c.JSON(200, gin.H{"ok": true})
// no envelope, no request_id
```

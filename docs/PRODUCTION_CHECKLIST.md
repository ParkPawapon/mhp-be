# Production Checklist

Use this as a baseline for STIN Smart Care. Adjust items based on your actual infra (VM, Docker, Kubernetes, or managed services).

## 1) Environment & Config
- APP_ENV=production
- HTTP_ENABLE_SWAGGER=false
- HTTP_ENABLE_METRICS=true
- CORS_ALLOWED_ORIGINS set to trusted web admin origins
- JWT_SECRET rotated and stored in secret manager
- SMS provider configured (SMS_PROVIDER=thaibulksms or disabled)
- TLS_CERT_FILE / TLS_KEY_FILE set when terminating TLS in-app
- OTEL_ENABLED and OTEL_SERVICE_NAME set if tracing is enabled

## 2) Secrets Management
- Never commit `.env`
- Use secret manager (GitHub Actions secrets, AWS SSM/Secrets, GCP Secret Manager, Vault)
- Rotate JWT secrets on incident or scheduled policy
- Restrict access to DB/Redis credentials

## 3) Database (PostgreSQL)
- Use managed Postgres or dedicated VM
- TLS enabled (DB_SSLMODE=require)
- Backups: daily full + point-in-time recovery
- Run migrations in CI/CD before deploy
- Verify required indexes exist
- Alert on slow queries and high connection usage

## 4) Cache (Redis)
- Use managed Redis or dedicated VM
- Enable AUTH and TLS if supported
- Configure eviction policy and persistence as needed
- Monitor memory usage and latency

## 5) Network & Security
- Enforce HTTPS at load balancer or app
- Restrict inbound ports to 443/80 (and 22 only if needed)
- Use WAF or rate limiting at edge when available
- Limit DB/Redis access to app network only

## 6) Logging & Observability
- Centralize logs (ELK, Loki, Datadog)
- Monitor `/metrics` with Prometheus + alerts
- Track 5xx rate, latency, and error codes
- Enable trace sampling if OpenTelemetry is used

## 7) CI/CD & Releases
- Build/test/lint in CI
- Tag releases (semver)
- Use blue/green or rolling deploys
- Validate health checks (`/healthz`, `/readyz`) before traffic switch

## 8) Backups & DR
- Document restore steps and RTO/RPO
- Test restore regularly
- Keep backup retention policy

## 9) Compliance & Data Safety
- Mask sensitive fields (citizen_id) for non-admin roles
- Never log OTP or passwords
- Apply audit log retention policy

## 10) Capacity & Scaling
- Set CPU/memory limits (K8s) or instance sizing (VM)
- Horizontal scaling with stateless API
- Ensure Redis and DB can handle peak load

## 11) Runbook
- Incident response checklist
- On-call contacts
- Deployment rollback steps

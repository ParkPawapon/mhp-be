# API CONTRACT

Base path: `/api/v1`

## Headers
- `Authorization: Bearer <token>` (required for protected endpoints)
- `X-Request-Id: <uuid>` (optional; generated if absent)

## Envelope
Success:
```json
{"data":{},"meta":{"request_id":"..."}}
```
Error:
```json
{"error":{"code":"...","message":"...","details":null},"meta":{"request_id":"..."}}
```
Pagination:
```json
{"data":[],"meta":{"request_id":"...","page":1,"page_size":20,"total":123}}
```

## Auth (Mobile)
### POST /auth/request-otp
Request:
```json
{"phone":"0812345678","purpose":"register"}
```
Response:
```json
{"data":{"ref_code":"AB1234","expires_at":"2026-01-20T12:00:00Z"},"meta":{"request_id":"..."}}
```

### POST /auth/verify-otp
Request:
```json
{"phone":"0812345678","ref_code":"AB1234","otp_code":"123456","purpose":"register"}
```
Response:
```json
{"data":{"verified":true},"meta":{"request_id":"..."}}
```

### POST /auth/register
Request:
```json
{"phone":"0812345678","ref_code":"AB1234","password":"***","first_name":"A","last_name":"B"}
```
Response:
```json
{"data":{"user_id":"uuid","access_token":"...","refresh_token":"..."},"meta":{"request_id":"..."}}
```

### POST /auth/login
Request:
```json
{"phone":"0812345678","password":"***"}
```
Response:
```json
{"data":{"access_token":"...","refresh_token":"..."},"meta":{"request_id":"..."}}
```

### POST /auth/forgot-password/request-otp
Request:
```json
{"phone":"0812345678"}
```
Response:
```json
{"data":{"ref_code":"AB1234","expires_at":"2026-01-20T12:00:00Z"},"meta":{"request_id":"..."}}
```

### POST /auth/forgot-password/confirm
Request:
```json
{"phone":"0812345678","ref_code":"AB1234","otp_code":"123456","new_password":"***"}
```
Response:
```json
{"data":{"updated":true},"meta":{"request_id":"..."}}
```

### POST /auth/refresh
Request:
```json
{"refresh_token":"..."}
```
Response:
```json
{"data":{"access_token":"...","refresh_token":"..."},"meta":{"request_id":"..."}}
```

### POST /auth/logout
Request:
```json
{"refresh_token":"..."}
```
Response:
```json
{"data":{"revoked":true},"meta":{"request_id":"..."}}
```

## User / Profile
### GET /me
Response:
```json
{"data":{"id":"uuid","role":"PATIENT","profile":{"first_name":"A","last_name":"B"}},"meta":{"request_id":"..."}}
```

### PATCH /me/profile
Request:
```json
{"first_name":"A","last_name":"B","hn":"HN001","citizen_id":"1234567890123","address_text":"..."}
```
Response:
```json
{"data":{"updated":true},"meta":{"request_id":"..."}}
```

### POST /me/device-tokens
Request:
```json
{"device_token":"...","platform":"ios"}
```
Response:
```json
{"data":{"saved":true},"meta":{"request_id":"..."}}
```

### PATCH /me/preferences
Request:
```json
{"weekly_reminder_enabled":true}
```
Response:
```json
{"data":{"weekly_reminder_enabled":true},"meta":{"request_id":"..."}}
```

## Support
### GET /support/emergency
Response:
```json
{"data":{"hotline":"1669","display_name":"Emergency 1669"},"meta":{"request_id":"..."}}
```

### POST /support/chat/requests
Request:
```json
{"message":"need help","category":"GENERAL","attachment_url":null}
```
Response:
```json
{"data":{"id":"uuid","status":"OPEN"},"meta":{"request_id":"..."}}
```

### GET /support/chat/requests
Response:
```json
{"data":[{"id":"uuid","user_id":"uuid","message":"need help","category":"GENERAL","status":"OPEN","created_at":"2026-01-20T12:00:00Z"}],"meta":{"request_id":"...","page":1,"page_size":20,"total":1}}
```

## Caregiver
### POST /caregivers/assignments
Request:
```json
{"patient_id":"uuid","caregiver_id":"uuid","relationship":"family"}
```
Response:
```json
{"data":{"assignment_id":"uuid"},"meta":{"request_id":"..."}}
```

### GET /caregivers/assignments?patient_id=
Response:
```json
{"data":[{"caregiver_id":"uuid","relationship":"family"}],"meta":{"request_id":"..."}}
```

## Medicines
### GET /medicines/categories
Response:
```json
{"data":[{"id":"uuid","name":"Hypertension","code":"HYPERTENSION","is_active":true}],"meta":{"request_id":"..."}}
```

### GET /medicines/categories/:id/items
Response:
```json
{"data":[{"id":"uuid","category_id":"uuid","display_name":"Amlodipine 5 mg","default_dosage_text":"1","is_active":true}],"meta":{"request_id":"..."}}
```

### GET /medicines/dosage-options
Response:
```json
{"data":["1/4","1/2","1","2"],"meta":{"request_id":"..."}}
```

### GET /medicines/meal-timing-options
Response:
```json
{"data":["BEFORE_MEAL","AFTER_MEAL","AFTER_MEAL_IMMEDIATELY","BEFORE_BED","UNTIL_FINISHED","NO_MILK","OTHER"],"meta":{"request_id":"..."}}
```

### GET /medicines/master
Response:
```json
{"data":[{"id":"uuid","trade_name":"..."}],"meta":{"request_id":"...","page":1,"page_size":20,"total":100}}
```

### POST /medicines/patient
Request:
```json
{"medicine_master_id":"uuid","category_item_id":"uuid","custom_name":"Amlodipine","dosage_amount":"1"}
```
Response:
```json
{"data":{"id":"uuid"},"meta":{"request_id":"..."}}
```

### GET /medicines/patient
Response:
```json
{"data":[{"id":"uuid","dosage_amount":"1 tab"}],"meta":{"request_id":"..."}}
```

### PATCH /medicines/patient/:id
Request:
```json
{"dosage_amount":"2 tabs"}
```
Response:
```json
{"data":{"updated":true},"meta":{"request_id":"..."}}
```

### DELETE /medicines/patient/:id
Response:
```json
{"data":{"deleted":true},"meta":{"request_id":"..."}}
```

### POST /medicines/patient/:id/schedules
Request:
```json
{"time_slot":"08:00","meal_timing":"BEFORE_MEAL"}
```
Response:
```json
{"data":{"schedule_id":"uuid"},"meta":{"request_id":"..."}}
```

### DELETE /medicines/schedules/:id
Response:
```json
{"data":{"deleted":true},"meta":{"request_id":"..."}}
```

## Notifications
### GET /notifications/upcoming?from=&to=
Response:
```json
{"data":[{"id":"uuid","template_code":"MED_BEFORE_MEAL_5MIN","scheduled_at":"2026-01-20T11:55:00Z","status":"PENDING"}],"meta":{"request_id":"..."}}
```

## Intake
### POST /intake
Request:
```json
{"schedule_id":"uuid","target_date":"2026-01-20","status":"TAKEN"}
```
Response:
```json
{"data":{"id":"uuid"},"meta":{"request_id":"..."}}
```

### GET /intake/history?from=&to=&user_id=
Response:
```json
{"data":[{"id":"uuid","status":"TAKEN"}],"meta":{"request_id":"..."}}
```

## Health Records & Assessments
### POST /health/records
Request:
```json
{"record_date":"2026-01-20","systolic_bp":120}
```
Response:
```json
{"data":{"id":"uuid"},"meta":{"request_id":"..."}}
```

### GET /health/records?from=&to=&user_id=
Response:
```json
{"data":[{"record_date":"2026-01-20"}],"meta":{"request_id":"..."}}
```

### POST /assessments/daily
Request:
```json
{"log_date":"2026-01-20","exercise_minutes":30}
```
Response:
```json
{"data":{"id":"uuid"},"meta":{"request_id":"..."}}
```

### GET /assessments/daily?from=&to=&user_id=
Response:
```json
{"data":[{"log_date":"2026-01-20"}],"meta":{"request_id":"..."}}
```

## Appointments
### GET /appointments?user_id=
Response:
```json
{"data":[{"id":"uuid","status":"PENDING"}],"meta":{"request_id":"..."}}
```

### POST /appointments
Request:
```json
{"title":"Checkup","appt_type":"HOSPITAL","appt_datetime":"2026-01-20T09:00:00Z"}
```
Response:
```json
{"data":{"id":"uuid"},"meta":{"request_id":"..."}}
```

### PATCH /appointments/:id/status
Request:
```json
{"status":"CONFIRMED"}
```
Response:
```json
{"data":{"updated":true},"meta":{"request_id":"..."}}
```

### DELETE /appointments/:id
Response:
```json
{"data":{"deleted":true},"meta":{"request_id":"..."}}
```

### POST /appointments/:id/notes
Request:
```json
{"visit_details":"..."}
```
Response:
```json
{"data":{"created":true},"meta":{"request_id":"..."}}
```

## Visits
### GET /visits/history?user_id=
Response:
```json
{"data":[{"appointment_id":"uuid","visit_note_id":"uuid","appt_datetime":"2026-01-20T12:00:00Z","title":"Home Visit","nurse_id":"uuid","visit_details":"...","created_at":"2026-01-20T12:30:00Z"}],"meta":{"request_id":"..."}}
```

## Health Content
### GET /content/health/categories
Response:
```json
{"data":["HYPERTENSION_KNOWLEDGE","HYPERTENSION_CONTROL","ABNORMAL_SYMPTOMS","CV_RISK_SCORE","FOOD_AND_MEDICINE","EXERCISE_AND_BMI","STRESS_MANAGEMENT","SLEEP"],"meta":{"request_id":"..."}}
```

### GET /content/health?published=true
Response:
```json
{"data":[{"id":"uuid","title":"..."}],"meta":{"request_id":"..."}}
```

### POST /content/health
Request:
```json
{"title":"...","body_content":"..."}
```
Response:
```json
{"data":{"id":"uuid"},"meta":{"request_id":"..."}}
```

### PATCH /content/health/:id
Request:
```json
{"title":"..."}
```
Response:
```json
{"data":{"updated":true},"meta":{"request_id":"..."}}
```

### POST /content/health/:id/publish
Request:
```json
{"is_published":true}
```
Response:
```json
{"data":{"published":true},"meta":{"request_id":"..."}}
```

## Admin/Staff (Web)
### POST /staff/login
Request:
```json
{"username":"admin","password":"***"}
```
Response:
```json
{"data":{"access_token":"...","refresh_token":"..."},"meta":{"request_id":"..."}}
```

### GET /admin/patients
Response:
```json
{"data":[{"id":"uuid","first_name":"A"}],"meta":{"request_id":"...","page":1,"page_size":20,"total":100}}
```

### GET /admin/patients/:id
Response:
```json
{"data":{"id":"uuid","first_name":"A"},"meta":{"request_id":"..."}}
```

### GET /admin/adherence?patient_id=&from=&to=
Response:
```json
{"data":[{"target_date":"2026-01-20","status":"TAKEN"}],"meta":{"request_id":"..."}}
```

## Audit
### GET /admin/audit-logs?from=&to=&actor_id=&action_type=
Response:
```json
{"data":[{"id":"uuid","action_type":"LOGIN"}],"meta":{"request_id":"...","page":1,"page_size":20,"total":100}}
```

## Health & Observability
### GET /healthz
Response:
```json
{"data":{"status":"ok"},"meta":{"request_id":"..."}}
```

### GET /readyz
Response:
```json
{"data":{"database":"ok","redis":"ok"},"meta":{"request_id":"..."}}
```

### GET /metrics
Response: Prometheus metrics text format.

## Error Examples
```json
{"error":{"code":"AUTH_UNAUTHORIZED","message":"unauthorized","details":null},"meta":{"request_id":"..."}}
```
```json
{"error":{"code":"VALIDATION_FAILED","message":"validation failed","details":{"phone":"required"}},"meta":{"request_id":"..."}}
```

## RBAC Matrix
| Endpoint Group | PATIENT | CAREGIVER | NURSE | ADMIN |
| --- | --- | --- | --- | --- |
| Auth | Yes | Yes | Yes | Yes |
| /me | Self | Self | Self | Self |
| /me/preferences | Self | Self | Self | Self |
| Caregiver assignments | No | No | Yes | Yes |
| Medicines/Intake | Self | Read assigned | Yes | Yes |
| Health records/assessments | Self | Read assigned | Yes | Yes |
| Appointments | Self | Read assigned | Yes | Yes |
| Visits history | Self | Read assigned | Yes | Yes |
| Health content | Read published | Read published | Create/Update | Full |
| Support emergency | Yes | Yes | Yes | Yes |
| Support chat | Create | No | List | List |
| Notifications | Self | Self | Self | Self |
| Admin endpoints | No | No | No | Yes |
| Audit logs | No | No | No | Yes |

## Sensitive Data Policy
- `password_hash` never returned.
- `citizen_id` masked for all non-admin roles.
- Caregiver can only access assigned patient data and only high-level fields.

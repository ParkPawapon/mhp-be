package constants

const (
	AuthUnauthorized       = "AUTH_UNAUTHORIZED"
	AuthForbidden          = "AUTH_FORBIDDEN"
	AuthInvalidCredentials = "AUTH_INVALID_CREDENTIALS"
	AuthOTPExpired         = "AUTH_OTP_EXPIRED"
	AuthOTPInvalid         = "AUTH_OTP_INVALID"
	AuthOTPUsed            = "AUTH_OTP_USED"
	AuthTokenInvalid       = "AUTH_TOKEN_INVALID"
	AuthTokenExpired       = "AUTH_TOKEN_EXPIRED"

	UserNotFound    = "USER_NOT_FOUND"
	UserConflict    = "USER_CONFLICT"
	UserInvalid     = "USER_INVALID"
	ProfileInvalid  = "USER_PROFILE_INVALID"

	MedInvalid   = "MED_INVALID"
	MedNotFound  = "MED_NOT_FOUND"

	ApptInvalid  = "APPT_INVALID"
	ApptNotFound = "APPT_NOT_FOUND"

	HealthInvalid  = "HEALTH_INVALID"
	HealthNotFound = "HEALTH_NOT_FOUND"

	ContentInvalid  = "CONTENT_INVALID"
	ContentNotFound = "CONTENT_NOT_FOUND"

	AuditInvalid = "AUDIT_INVALID"

	RateLimited = "RATE_LIMITED"

	ValidationFailed = "VALIDATION_FAILED"

	InternalError       = "INTERNAL_ERROR"
	InternalUnavailable = "INTERNAL_UNAVAILABLE"
	InternalNotImplemented = "INTERNAL_NOT_IMPLEMENTED"
)

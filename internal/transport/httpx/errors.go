package httpx

import (
	"net/http"
	"strings"

	"github.com/ParkPawapon/mhp-be/internal/constants"
	"github.com/ParkPawapon/mhp-be/internal/domain"
)

type HTTPError struct {
	Status  int
	Code    string
	Message string
	Details any
}

func MapError(err error) HTTPError {
	appErr, ok := domain.AsAppError(err)
	if !ok {
		return HTTPError{
			Status:  http.StatusInternalServerError,
			Code:    constants.InternalError,
			Message: "internal error",
			Details: nil,
		}
	}

	status := statusFromCode(appErr.Code)
	message := appErr.Message
	if message == "" {
		message = "error"
	}

	return HTTPError{
		Status:  status,
		Code:    appErr.Code,
		Message: message,
		Details: appErr.Details,
	}
}

func statusFromCode(code string) int {
	switch code {
	case constants.AuthForbidden:
		return http.StatusForbidden
	case constants.AuthUnauthorized, constants.AuthInvalidCredentials, constants.AuthTokenInvalid, constants.AuthTokenExpired:
		return http.StatusUnauthorized
	case constants.AuthOTPExpired, constants.AuthOTPInvalid, constants.AuthOTPUsed:
		return http.StatusBadRequest
	case constants.UserConflict:
		return http.StatusConflict
	case constants.RateLimited:
		return http.StatusTooManyRequests
	case constants.ValidationFailed, constants.MedInvalid, constants.ApptInvalid, constants.HealthInvalid, constants.ContentInvalid:
		return http.StatusBadRequest
	case constants.UserNotFound, constants.MedNotFound, constants.ApptNotFound, constants.HealthNotFound, constants.ContentNotFound:
		return http.StatusNotFound
	case constants.InternalNotImplemented:
		return http.StatusNotImplemented
	case constants.InternalUnavailable:
		return http.StatusServiceUnavailable
	default:
		switch {
		case strings.HasPrefix(code, "AUTH_"):
			return http.StatusUnauthorized
		case strings.HasPrefix(code, "VALIDATION_"):
			return http.StatusBadRequest
		case strings.HasPrefix(code, "RATE_"):
			return http.StatusTooManyRequests
		case strings.HasPrefix(code, "INTERNAL_"):
			return http.StatusInternalServerError
		default:
			return http.StatusBadRequest
		}
	}
}

package domain

import "errors"

type AppError struct {
	Code    string
	Message string
	Err     error
	Details any
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func NewError(code, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

func WrapError(code, message string, err error) *AppError {
	return &AppError{Code: code, Message: message, Err: err}
}

func WithDetails(err *AppError, details any) *AppError {
	err.Details = details
	return err
}

func AsAppError(err error) (*AppError, bool) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}

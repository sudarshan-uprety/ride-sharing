package errors

import (
	"fmt"
	"net/http"
)

type ErrorType string

const (
	ErrorTypeValidation   ErrorType = "VALIDATION_ERROR"
	ErrorTypeConflict     ErrorType = "CONFLICT_ERROR"
	ErrorTypeNotFound     ErrorType = "NOT_FOUND_ERROR"
	ErrorTypeUnauthorized ErrorType = "UNAUTHORIZED_ERROR"
	ErrorTypeForbidden    ErrorType = "FORBIDDEN_ERROR"
	ErrorTypeInternal     ErrorType = "INTERNAL_ERROR"
)

type AppError struct {
	Type    ErrorType `json:"type"`
	Message string    `json:"message"`
	Details any       `json:"details,omitempty"`
	Err     error     `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Type, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func NewValidationError(message string, details any) *AppError {
	return &AppError{
		Type:    ErrorTypeValidation,
		Message: message,
		Details: details,
	}
}

func NewConflictError(message string) *AppError {
	return &AppError{
		Type:    ErrorTypeConflict,
		Message: message,
	}
}

func NewNotFoundError(message string) *AppError {
	return &AppError{
		Type:    ErrorTypeNotFound,
		Message: message,
	}
}

func NewUnauthorizedError(message string) *AppError {
	return &AppError{
		Type:    ErrorTypeUnauthorized,
		Message: message,
	}
}

func NewForbiddenError(message string) *AppError {
	return &AppError{
		Type:    ErrorTypeForbidden,
		Message: message,
	}
}

func NewInternalError(err error) *AppError {
	return &AppError{
		Type:    ErrorTypeInternal,
		Message: "internal server error",
		Err:     err,
	}
}

func HTTPStatusFromErrorType(t ErrorType) int {
	switch t {
	case ErrorTypeValidation:
		return http.StatusBadRequest
	case ErrorTypeConflict:
		return http.StatusConflict
	case ErrorTypeNotFound:
		return http.StatusNotFound
	case ErrorTypeUnauthorized:
		return http.StatusUnauthorized
	case ErrorTypeForbidden:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}

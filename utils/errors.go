package utils

// import (
// 	"net/http"
// )

// // ValidationErrorType represents the type of validation error
// type ValidationErrorType struct {
// 	StatusCode int
// 	Message    string
// 	Details    interface{}
// }

// // Common validation error types
// var (
// 	ErrBadRequest = func(details interface{}) ValidationErrorType {
// 		return ValidationErrorType{http.StatusBadRequest, "Validation failed", details}
// 	}
// 	ErrResourceExists = func(details interface{}) ValidationErrorType {
// 		return ValidationErrorType{http.StatusConflict, "Resource already exists", details}
// 	}
// 	ErrPasswordMismatch = func(details interface{}) ValidationErrorType {
// 		return ValidationErrorType{http.StatusBadRequest, "Passwords do not match", details}
// 	}
// 	ErrDatabaseError = func(details interface{}) ValidationErrorType {
// 		return ValidationErrorType{http.StatusInternalServerError, "Internal server error", details}
// 	}
// 	ErrServerError = func(details interface{}) ValidationErrorType {
// 		return ValidationErrorType{http.StatusInternalServerError, "Server error", details}
// 	}
// )

// // HandleValidationError is a common function to handle validation errors and return appropriate responses
// func HandleValidationError(w http.ResponseWriter, errType ValidationErrorType) {
// 	Error(w, errType.StatusCode, errType.Message, errType.Details)
// }

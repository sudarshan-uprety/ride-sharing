package utils

import "net/http"

// HTTP Status Codes
const (
	SUCCESS_CODE                       = http.StatusOK
	SUCCESS_FETCH_CODE                 = http.StatusOK
	SUCCESS_UPDATED_CODE               = http.StatusOK
	SUCCESS_CREATED_CODE               = http.StatusCreated
	SUCCESS_DELETED_CODE               = http.StatusNoContent
	ERROR_BAD_REQUEST_CODE             = http.StatusBadRequest
	ERROR_UNAUTHORIZED_CODE            = http.StatusUnauthorized
	ERROR_FORBIDDEN_CODE               = http.StatusForbidden
	ERROR_NOT_FOUND                    = http.StatusNotFound
	ERROR_RESOURCE_ALREADY_EXISTS_CODE = http.StatusConflict
	ERROR_INTERNAL_CODE                = http.StatusInternalServerError
	SERVICE_UNAVAILABLE_CODE           = http.StatusServiceUnavailable
	UNPROCESSABLE_ENTITY_CODE          = http.StatusUnprocessableEntity
)

// Success Messages
// Success Messages
const (
	AUTH_SUCCESS_LOGIN          = "Login successful."
	AUTH_SUCCESS_LOGOUT         = "Logout successful."
	AUTH_SUCCESS_REGISTER       = "Registration successful."
	AUTH_SUCCESS_PASSWORD_RESET = "Password reset successful."
	AUTH_SUCCESS_EMAIL_SENT     = "Verification email sent successfully."
	AUTH_SUCCESS_VERIFIED       = "Account verified successfully."
	AUTH_SUCCESS_TOKEN_REFRESH  = "Token refreshed successfully."
	AUTH_SUCCESS_UPDATE_PROFILE = "Profile updated successfully."
)

// Error Messages
const (
	ERROR_DOES_NOT_EXIST                  = "Resource Does Not Exist."
	ERROR_RESOURCE_ALREADY_EXISTS         = "Resource Already Exist"
	ERROR_SERVER_DOWN                     = "Server Down"
	INTERNAL_SERVER_ERROR                 = "Internal Server Error"
	EMAIL_SERVER_ERROR                    = "Email Server Error"
	ErrInvalidInput                       = "Invalid input"
	ErrEmailAlreadyUsed                   = "Email already used"
	ErrHashingPassword                    = "Failed to hash password"
	ErrUserCreate                         = "Failed to create user"
	DATABASE_NOT_AVAILABLE_FOR_CONNECTION = "Internal Server Error"
)

// Database connection error phrases
var DATABASE_CONNECTION_ERRORS = []string{
	"is the server running",
	"failure in name resolution",
	"connection refused",
}

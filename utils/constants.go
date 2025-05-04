package utils

// HTTP Status Codes
const (
	SUCCESS_CODE                       = 200
	SUCCESS_FETCH_CODE                 = 200
	SUCCESS_UPDATED_CODE               = 200
	SUCCESS_CREATED_CODE               = 201
	SUCCESS_DELETED_CODE               = 204
	ERROR_BAD_REQUEST_CODE             = 400
	ERROR_UNAUTHORIZED_CODE            = 401
	ERROR_FORBIDDEN_CODE               = 403
	ERROR_NOT_FOUND                    = 404
	ERROR_RESOURCE_ALREADY_EXISTS_CODE = 409
	ERROR_FOUND_CODE                   = 302
	ERROR_INTERNAL_CODE                = 500
	SERVICE_UNAVAILABLE_CODE           = 503
	UNPROCESSABLE_ENTITY_CODE          = 422
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
	ERROR_DOES_NOT_EXIST          = "Resource Does Not Exist."
	ERROR_RESOURCE_ALREADY_EXISTS = "Resource Already Exist"
	ERROR_SERVER_DOWN             = "Server Down"
	INTERNAL_SERVER_ERROR         = "Internal Server Error"
	EMAIL_SERVER_ERROR            = "Email Server Error"
	ErrInvalidInput               = "Invalid input"
	ErrEmailAlreadyUsed           = "Email already used"
	ErrHashingPassword            = "Failed to hash password"
	ErrUserCreate                 = "Failed to create user"
)

// Database connection error phrases
var DATABASE_CONNECTION_ERRORS = []string{
	"is the server running",
	"failure in name resolution",
	"connection refused",
}

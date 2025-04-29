package constants

// HTTP Status Codes
const (
	SUCCESS              = 200
	SUCCESS_FETCH        = 200
	SUCCESS_UPDATED      = 200
	SUCCESS_CREATED      = 201
	SUCCESS_DELETED      = 204
	ERROR_BAD_REQUEST    = 400
	ERROR_UNAUTHORIZED   = 401
	ERROR_FORBIDDEN      = 403
	ERROR_NOT_FOUND      = 404
	ERROR_FOUND          = 302
	ERROR_INTERNAL       = 500
	SERVICE_UNAVAILABLE  = 503
	UNPROCESSABLE_ENTITY = 422
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
	ERROR_DOES_NOT_EXIST           = "Resource Does Not Exist."
	ERROR_SERVER_DOWN              = "Server Down"
	INTERNAL_SERVER_ERROR          = "Internal Server Error"
	EMAIL_SERVER_ERROR             = "Email Server Error"
	DATABASE_NOT_AVAILABLE_MESSAGE = "Sorry, but the database is either offline or not accepting connections."
)

// Database connection error phrases
var DATABASE_CONNECTION_ERRORS = []string{
	"is the server running",
	"failure in name resolution",
	"connection refused",
}

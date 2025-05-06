package response

import (
	"ride-sharing/internal/pkg/errors"
	"ride-sharing/internal/pkg/validation"

	"github.com/gin-gonic/gin"
)

type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Meta    interface{} `json:"meta,omitempty"`
}

type ErrorResponse struct {
	Success bool        `json:"success"`
	Error   string      `json:"error"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

func Success(c *gin.Context, status int, message string, data interface{}, meta interface{}) {
	c.JSON(status, SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

func Error(c *gin.Context, appErr *errors.AppError) {
	statusCode := errors.HTTPStatusFromErrorType(appErr.Type)

	// Special handling for validation errors to provide better details
	if appErr.Type == errors.ErrorTypeValidation {
		details := make(map[string]string)

		// If the error details are already a map[string]string, use them directly
		if detailsMap, ok := appErr.Details.(map[string]string); ok {
			details = detailsMap
		} else if validationErr, ok := appErr.Details.(error); ok {
			// If the details is an error, try to process it as a validation error
			details = validation.ProcessValidationError(validationErr)
		} else if detailsStr, ok := appErr.Details.(string); ok && appErr.Details != nil {
			// If it's just a string, add it as a general error
			details["_error"] = detailsStr
		}

		c.JSON(statusCode, gin.H{
			"success": false,
			"error":   string(appErr.Type),
			"message": "Validation failed",
			"details": details,
		})
		return
	}

	// For other error types
	c.JSON(statusCode, gin.H{
		"success": false,
		"error":   string(appErr.Type),
		"message": appErr.Message,
		"details": appErr.Details,
	})
}

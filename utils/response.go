package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorResponse represents a standardized error response
type ErrorResponseStruct struct {
	StatusCode int         `json:"-"`
	Message    string      `json:"message"`
	Details    string      `json:"details,omitempty"`
	Errors     interface{} `json:"errors,omitempty"`
}

func (e ErrorResponseStruct) Error() string {
	return e.Message
}

func NewErrorResponse(statusCode int, message, details string, errors interface{}) *ErrorResponseStruct {
	return &ErrorResponseStruct{
		StatusCode: statusCode,
		Message:    message,
		Details:    details,
		Errors:     errors,
	}
}

// SuccessResponse sends a standardized success response
func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}, warning string) {
	c.JSON(statusCode, gin.H{
		"success": true,
		"message": message,
		"data":    data,
		"warning": warning,
	})
}

// ErrorResponse sends a standardized error response
func ErrorResponse(c *gin.Context, err error) {
	// Default error response
	statusCode := http.StatusBadRequest
	message := "Validation failed"
	// details := err.Error()
	errors := make(map[string]string)

	// If it's our custom error type, use its properties
	if errResp, ok := err.(*ErrorResponseStruct); ok {
		statusCode = errResp.StatusCode
		message = errResp.Message
		if errResp.Errors != nil {
			errors = errResp.Errors.(map[string]string)
		}
	}

	c.JSON(statusCode, gin.H{
		"success": false,
		"message": message,
		"errors":  errors,
	})
}

package middleware

import (
	"net/http"
	"ride-sharing/internal/pkg/errors"
	"ride-sharing/internal/pkg/validation"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// ValidationMiddleware intercepts validation errors from binding and provides a better formatted response
func ValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there were any validation errors set during the request
		if len(c.Errors) > 0 {
			// Look for validation errors
			for _, err := range c.Errors {
				// Check if it's a validation error
				if valErrs, ok := err.Err.(validator.ValidationErrors); ok {
					details := validation.ProcessValidationError(valErrs)
					c.JSON(http.StatusBadRequest, gin.H{
						"success": false,
						"error":   string(errors.ErrorTypeValidation),
						"message": "Validation failed",
						"details": details,
					})
					c.Abort()
					return
				}
			}
		}
	}
}

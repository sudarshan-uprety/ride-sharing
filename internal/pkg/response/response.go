package response

import (
	"ride-sharing/internal/pkg/errors"

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

func Error(c *gin.Context, err *errors.AppError) {
	status := errors.HTTPStatusFromErrorType(err.Type)
	c.JSON(status, ErrorResponse{
		Success: false,
		Error:   string(err.Type),
		Message: err.Message,
		Details: err.Details,
	})
}

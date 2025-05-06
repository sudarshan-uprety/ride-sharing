package middleware

import (
	"context"
	"time"

	"ride-sharing/internal/pkg/logging"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// LoggingMiddleware is a gin middleware that logs request details
// and adds request_id and correlation_id to both context and response headers
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Get or generate tracking IDs
		requestID := c.GetHeader(logging.RequestIDKey)
		if requestID == "" {
			requestID = uuid.New().String()
		}

		correlationID := c.GetHeader(logging.CorrelationID)
		if correlationID == "" {
			correlationID = uuid.New().String()
		}

		// Add IDs to context and headers
		ctx := context.WithValue(c.Request.Context(), logging.RequestIDKey, requestID)
		ctx = context.WithValue(ctx, logging.CorrelationID, correlationID)
		c.Request = c.Request.WithContext(ctx)

		// Set response headers for tracking
		c.Writer.Header().Set(logging.RequestIDKey, requestID)
		c.Writer.Header().Set(logging.CorrelationID, correlationID)

		// Get logger with request context
		logger := logging.GetLogger().WithContext(ctx)

		// Process request
		c.Next()

		// Collect metrics after request completes
		latency := time.Since(start)
		status := c.Writer.Status()
		bodySize := c.Writer.Size()
		if bodySize < 0 {
			bodySize = 0
		}

		// Prepare standard log fields
		fields := []zap.Field{
			zap.Int("status", status),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("query", c.Request.URL.RawQuery),
			zap.String("ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
			zap.Duration("latency", latency),
			zap.Int("body_size", bodySize),
		}

		// Log based on response status code
		switch {
		case status >= 400 && status < 500:
			if len(c.Errors) > 0 {
				fields = append(fields, zap.String("error", c.Errors.String()))
			}
			logger.Warn("client error", fields...)
		case status >= 500:
			fields = append(fields, zap.String("error", c.Errors.String()))
			logger.Error("server error", fields...)
		default:
			logger.Info("request processed", fields...)
		}
	}
}

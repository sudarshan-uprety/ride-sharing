package middleware

import (
	"context"
	"time"

	"ride-sharing/internal/pkg/logging"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func GinLoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Get or generate IDs
		requestID := c.GetHeader(logging.RequestIDKey)
		if requestID == "" {
			requestID = uuid.New().String()
		}

		correlationID := c.GetHeader(logging.CorrelationID)
		if correlationID == "" {
			correlationID = uuid.New().String()
		}

		// Add to context and headers
		ctx := context.WithValue(c.Request.Context(), logging.RequestIDKey, requestID)
		ctx = context.WithValue(ctx, logging.CorrelationID, correlationID)
		c.Request = c.Request.WithContext(ctx)

		c.Writer.Header().Set(logging.RequestIDKey, requestID)
		c.Writer.Header().Set(logging.CorrelationID, correlationID)

		// Process request
		c.Next()

		// Collect metrics after request completes
		latency := time.Since(start)
		status := c.Writer.Status()
		bodySize := c.Writer.Size()
		if bodySize < 0 {
			bodySize = 0
		}

		logger := logging.GetLogger().WithContext(ctx)

		// Single log entry per request with all details
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

		switch {
		case status >= 400 && status < 500:
			logger.Warn("client error", fields...)
		case status >= 500:
			fields = append(fields, zap.String("error", c.Errors.String()))
			logger.Error("server error", fields...)
		default:
			logger.Info("request processed", fields...)
		}
	}
}

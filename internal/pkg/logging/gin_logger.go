package logging

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GinLogger returns a gin.HandlerFunc that logs requests
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Collect metrics after request completes
		latency := time.Since(start)
		status := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		userAgent := c.Request.UserAgent()
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		// Get logger from context or create new
		logger := GetLogger().WithContext(c.Request.Context())

		// Log based on status code
		switch {
		case status >= 400 && status < 500:
			logger.Warn("client error",
				zap.Int("status", status),
				zap.String("method", method),
				zap.String("path", path),
				zap.String("query", query),
				zap.String("ip", clientIP),
				zap.String("user_agent", userAgent),
				zap.Duration("latency", latency),
				zap.String("error", errorMessage),
			)
		case status >= 500:
			logger.Error("server error",
				zap.Int("status", status),
				zap.String("method", method),
				zap.String("path", path),
				zap.String("query", query),
				zap.String("ip", clientIP),
				zap.String("user_agent", userAgent),
				zap.Duration("latency", latency),
				zap.String("error", errorMessage),
			)
		default:
			logger.Info("request",
				zap.Int("status", status),
				zap.String("method", method),
				zap.String("path", path),
				zap.String("query", query),
				zap.String("ip", clientIP),
				zap.String("user_agent", userAgent),
				zap.Duration("latency", latency),
			)
		}
	}
}

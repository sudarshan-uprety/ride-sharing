package middleware

import (
	"ride-sharing/internal/pkg/auth"

	"github.com/gin-gonic/gin"
)

func RiderOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractToken(c)
		if token == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized: Missing token"})
			return
		}

		claims, err := auth.ValidateToken(token, "RIDER_JWT_SECRET")
		if err != nil {
			c.AbortWithStatusJSON(403, gin.H{"error": "Forbidden: Invalid token"})
			return
		}

		c.Set("riderID", claims.UserID)
		c.Next()
	}
}

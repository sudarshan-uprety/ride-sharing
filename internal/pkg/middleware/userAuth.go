// middleware/auth.go
package middleware

import (
	"ride-sharing/internal/pkg/auth"

	"github.com/gin-gonic/gin"
)

func UserOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractToken(c)
		if token == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized: Missing token"})
			return
		}

		claims, err := auth.ValidateToken(token, "USER_JWT_SECRET")
		if err != nil {
			c.AbortWithStatusJSON(403, gin.H{"error": "Forbidden: Invalid token"})
			return
		}

		c.Set("userID", claims.UserID)
		c.Next()
	}
}

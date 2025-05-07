// middleware/auth_utils.go
package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func extractToken(c *gin.Context) string {
	bearerToken := c.GetHeader("Authorization")
	return strings.TrimPrefix(bearerToken, "Bearer ")
}

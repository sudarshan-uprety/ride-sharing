package auth

import (
	"fmt"
	"log"
	"ride-sharing/internal/pkg/errors"
	"ride-sharing/internal/pkg/response"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID string `json:"sub"`
	jwt.RegisteredClaims
}

func ValidateToken(tokenString, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			response.Error(c, errors.NewUnauthorizedError("missing or invalid authorization header"))
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := ValidateToken(tokenStr, secret)
		if err != nil {
			log.Println("ERROR IS-----------------------------------------------", err)
			response.Error(c, errors.NewUnauthorizedError("invalid or expired token"))
			c.Abort()
			return
		}

		// Store user ID in context
		c.Set("userID", claims.UserID)
		c.Next()
	}
}

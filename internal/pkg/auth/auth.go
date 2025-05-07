package auth

import (
	"fmt"

	userRepository "ride-sharing/internal/domains/users/repository"
	"ride-sharing/internal/pkg/errors"
	"ride-sharing/internal/pkg/response"

	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID    string   `json:"sub"`
	TokenType string   `json:"typ"`
	UserRole  UserType `json:"role"`
	jwt.RegisteredClaims
}

func ValidateToken(tokenStr string, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
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
			response.Error(c, errors.NewUnauthorizedError("invalid or expired token"))
			c.Abort()
			return
		}

		if claims.TokenType != "access" {
			response.Error(c, errors.NewUnauthorizedError("invalid token type: access token required"))
			c.Abort()
			return
		}

		// Fetch complete user based on type
		userRole := UserType(claims.UserRole)
		if !userRole.IsValid() {
			response.Error(c, errors.NewUnauthorizedError("invalid user role"))
			c.Abort()
			return
		}

		var user interface{}
		switch userRole {
		case UserTypeUser:
			user, err = userRepository.UserRepository.GetByID(c.Request.Context(), claims.UserID)
		// case UserTypeRider:
		//     user, err = riderRepo.GetByID(c.Request.Context(), claims.UserID)
		// case UserTypeAdmin:
		//     user, err = adminRepo.GetByID(c.Request.Context(), claims.UserID)
		default:
			response.Error(c, errors.NewUnauthorizedError("unsupported user role"))
			c.Abort()
			return
		}

		if err != nil {
			response.Error(c, errors.NewInternalError(err))
			c.Abort()
			return
		}

		if user == nil {
			response.Error(c, errors.NewUnauthorizedError("user not found"))
			c.Abort()
			return
		}

		// Store in context
		c.Set("userID", claims.UserID)
		c.Set("role", userRole)
		c.Set("authUser", user) // Store the concrete user object
		c.Next()
	}
}

package middleware

import (
	"fmt"
	"strings"
	"time"

	"ride-sharing/internal/domains/users/models"
	"ride-sharing/internal/pkg/auth"
	"ride-sharing/internal/pkg/errors"
	"ride-sharing/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	tokenService  *auth.TokenService
	userProviders map[auth.UserType]auth.UserProvider
}

func NewAuthMiddleware(
	tokenService *auth.TokenService,
	userProviders map[auth.UserType]auth.UserProvider,
) *AuthMiddleware {
	return &AuthMiddleware{
		tokenService:  tokenService,
		userProviders: userProviders,
	}
}

func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, errors.NewUnauthorizedError("authorization header required"))
			c.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			response.Error(c, errors.NewUnauthorizedError("invalid authorization header format"))
			c.Abort()
			return
		}

		token := tokenParts[1]
		claims, err := m.tokenService.ValidateAccessToken(token)
		if err != nil {
			response.Error(c, errors.NewUnauthorizedError("invalid token: "+err.Error()))
			c.Abort()
			return
		}
		provider, exists := m.userProviders[claims.UserType]
		if !exists {
			response.Error(c, errors.NewUnauthorizedError("invalid user type"))
			c.Abort()
			return
		}

		user, err := provider.GetByID(c.Request.Context(), claims.UserID)
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

		userData, ok := user.(*models.User)
		if !ok {
			response.Error(c, errors.NewInternalError(fmt.Errorf("user is not of expected type *models.User")))
			c.Abort()
			return
		}

		// Parse claims.PasswordChangedAt (from token) to int64 (Unix timestamp in nanoseconds)
		tokenPasswordChangedAt := time.Unix(0, claims.PasswordChangedAt)

		// Compare PasswordChangedAt values
		if tokenPasswordChangedAt.Before(*userData.PasswordChangedAt) {
			response.Error(c, errors.NewUnauthorizedError("token is no longer valid due to password change"))
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("userType", claims.UserType)
		c.Set("authUser", user)
		c.Next()
	}
}

func RequireUserType(userType auth.UserType) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentType, exists := c.Get("userType")
		if !exists || currentType != userType {
			response.Error(c, errors.NewForbiddenError("access forbidden"))
			c.Abort()
			return
		}
		c.Next()
	}
}

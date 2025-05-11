package routes

import (
	"ride-sharing/internal/domains/users/delivery/http"
	"ride-sharing/internal/domains/users/repository"
	"ride-sharing/internal/domains/users/service"
	"ride-sharing/internal/pkg/auth"
	"ride-sharing/internal/pkg/middleware"
	"ride-sharing/internal/pkg/provider"
	"ride-sharing/internal/pkg/redis"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, tokenService *auth.TokenService, otpStore *redis.OTPStore) *gin.Engine {
	router := gin.Default()

	router.Use(middleware.LoggingMiddleware(), gin.Recovery())

	// Initialize dependencies
	userRepo := repository.NewUserRepository(db)
	// Create user providers
	userProviders := map[auth.UserType]auth.UserProvider{
		auth.UserTypeUser: provider.NewUserProvider(userRepo),
	}
	userService := service.NewUserService(userRepo, tokenService, otpStore, userProviders)
	userHandler := http.NewUserHandler(userService)

	authMiddleware := middleware.NewAuthMiddleware(tokenService, userProviders)

	// API versioning
	api := router.Group("/api/v1")

	// Public user routes
	userRoutes := api.Group("/users")
	{
		userRoutes.POST("/register", userHandler.Register)
		userRoutes.POST("/login", userHandler.Login)
		userRoutes.POST("/refresh", userHandler.Refresh)
		userRoutes.POST("/forget-password", userHandler.ForgetPassword)
		userRoutes.POST("/verify-reset", userHandler.VerifyForgetPassword)

	}

	// Protected user routes
	authRoutes := api.Group("/users")
	authRoutes.Use(authMiddleware.Authenticate(), middleware.RequireUserType(auth.UserTypeUser))
	{
		authRoutes.POST("/change-password", userHandler.ChangePassword)
	}

	return router
}

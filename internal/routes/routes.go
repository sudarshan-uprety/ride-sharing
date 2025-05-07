package routes

import (
	"ride-sharing/internal/domains/users/delivery/http"
	"ride-sharing/internal/domains/users/repository"
	"ride-sharing/internal/domains/users/service"
	"ride-sharing/internal/pkg/auth"
	"ride-sharing/internal/pkg/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, tokenService *auth.TokenService) *gin.Engine {
	router := gin.Default()

	router.Use(middleware.LoggingMiddleware(), gin.Recovery())

	// Initialize dependencies
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, tokenService)
	userHandler := http.NewUserHandler(userService)

	// API versioning
	api := router.Group("/api/v1")

	// Public user routes
	userRoutes := api.Group("/users")
	{
		userRoutes.POST("/register", userHandler.Register)
		userRoutes.POST("/login", userHandler.Login)
	}

	// Protected user routes
	authRoutes := api.Group("/users")
	authRoutes.Use(auth.AuthMiddleware("72bHG8VL0fRxXjsrfBx6o1Esz0Io0Kdb"))
	{
		authRoutes.POST("/change-password", userHandler.ChangePassword)
	}

	return router
}

package routes

import (
	"ride-sharing/internal/domains/users/delivery/http"
	"ride-sharing/internal/domains/users/repository"
	"ride-sharing/internal/domains/users/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	// Initialize dependencies
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := http.NewUserHandler(userService)

	// API versioning
	api := router.Group("/api/v1")

	// User routes
	userRoutes := api.Group("/users")
	{
		userRoutes.POST("/register", userHandler.Register)
		// userRoutes.POST("/login", userHandler.Login)
	}

	return router
}

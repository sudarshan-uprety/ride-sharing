package users

import (
	"ride-sharing/schemas/users"
	"ride-sharing/src/users/queries"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupUserRoutes(router *gin.Engine, db *gorm.DB) {
	// Initialize dependencies
	userRepo := queries.NewUserRepository(db)
	userValidator := users.NewUserValidator(userRepo) // Note: using users package
	userHandler := NewUserHandler(userRepo, userValidator)

	// Setup routes
	router.POST("/auth/signup", userHandler.RegisterUser)
}

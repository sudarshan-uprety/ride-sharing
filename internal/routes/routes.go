package routes

import (
	"ride-sharing/config"
	"ride-sharing/internal/domains/users/delivery/http"
	"ride-sharing/internal/domains/users/repository"
	"ride-sharing/internal/domains/users/service"
	"ride-sharing/internal/pkg/auth"
	email "ride-sharing/internal/pkg/grpcclient"
	"ride-sharing/internal/pkg/middleware"
	"ride-sharing/internal/pkg/provider"
	"ride-sharing/internal/pkg/redis"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, tokenService *auth.TokenService, otpStore *redis.OTPStore, notificationService *email.NotificationClient, cfg *config.Config) *gin.Engine {
	router := gin.Default()
	router.Use(middleware.LoggingMiddleware(), gin.Recovery())

	if cfg.Server.Environment != "production" {
		// Create dynamic Swagger handler
		swaggerHandler := ginSwagger.WrapHandler(
			swaggerFiles.Handler,
			ginSwagger.URL(cfg.Server.SwaggerURL+"/swagger/doc.json"),
			ginSwagger.DefaultModelsExpandDepth(-1),
		)
		router.GET("/swagger/*any", swaggerHandler)
	}
	// Initialize dependencies
	userRepo := repository.NewUserRepository(db)
	// Create user providers
	userProviders := map[auth.UserType]auth.UserProvider{
		auth.UserTypeUser: provider.NewUserProvider(userRepo),
	}
	userService := service.NewUserService(userRepo, tokenService, otpStore, notificationService, userProviders)
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
		userRoutes.POST("/verify-email", userHandler.VerifyEmail)

	}

	// Protected user routes
	authRoutes := api.Group("/users")
	authRoutes.Use(authMiddleware.Authenticate(), middleware.RequireUserType(auth.UserTypeUser))
	{
		authRoutes.POST("/change-password", userHandler.ChangePassword)
		authRoutes.GET("/profile", userHandler.UserProfile)
	}

	return router
}

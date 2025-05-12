package main

import (
	"log"
	"ride-sharing/config"
	_ "ride-sharing/docs"
	"ride-sharing/internal/domains/users/models"
	"ride-sharing/internal/pkg/auth"
	"ride-sharing/internal/pkg/database"
	"ride-sharing/internal/pkg/logging"
	"ride-sharing/internal/pkg/redis"
	"ride-sharing/internal/pkg/validation"
	"ride-sharing/internal/routes"
	"time"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// @title           Ride Sharing Auth API
// @version         1.0
// @description     This is a ride sharing service API.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
// @description                 Type "Bearer" followed by a space and JWT token
func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	logging.InitLogger(logging.LogConfig{
		Environment: cfg.Log.Environment,
		Version:     cfg.Log.Version,
		ServiceName: cfg.Log.ServiceName,
	})

	// Initialize database
	db, err := database.NewPostgresDB(database.DBConfig{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		User:     cfg.DB.User,
		Password: cfg.DB.Password,
		Name:     cfg.DB.Name,
	})
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	// Initialize Redis
	redisClient := redis.New(cfg)
	defer redisClient.Close()

	otpStore := redis.NewOTPStore(redisClient)
	// Initialize token service
	tokenService := auth.NewTokenService(
		cfg.JWT.AccessSecret,
		cfg.JWT.RefreshSecret,
		time.Hour*6,    // Access token expires in 1 hour
		time.Hour*24*7, // Refresh token expires in 1 week
	)

	// Auto-migrate models
	if err := database.AutoMigrate(db, &models.User{}); err != nil {
		log.Fatalf("failed to auto-migrate models: %v", err)
	}

	// Setup router
	router := routes.SetupRouter(db, tokenService, otpStore, cfg)

	// Register custom validators
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		validation.RegisterCustomValidators(v)
	}

	// Start server
	log.Printf("server starting on port %s", cfg.Server.Port)
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

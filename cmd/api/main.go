package main

import (
	"log"
	"ride-sharing/config"
	"ride-sharing/internal/domains/users/models"
	"ride-sharing/internal/pkg/database"
	"ride-sharing/internal/pkg/logging"
	"ride-sharing/internal/pkg/validation"
	"ride-sharing/internal/routes"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

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

	// Auto-migrate models
	if err := database.AutoMigrate(db, &models.User{}); err != nil {
		log.Fatalf("failed to auto-migrate models: %v", err)
	}

	// Setup router
	router := routes.SetupRouter(db)

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

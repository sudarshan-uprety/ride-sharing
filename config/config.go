package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
	}
	Server struct {
		Port string
	}
	JWT struct {
		AccessSecret  string
		RefreshSecret string
	}
	Log struct {
		Environment string
		Version     string
		ServiceName string
	}
}

func Load() (*Config, error) {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found - using system environment variables")
	}

	cfg := &Config{}

	// Database configuration
	cfg.DB.Host = getEnv("DB_HOST", "localhost")
	cfg.DB.Port = getEnv("DB_PORT", "5432")
	cfg.DB.User = getEnv("DB_USER", "postgres") // Default to postgres if not set
	cfg.DB.Password = getEnv("DB_PASSWORD", "postgres")
	cfg.DB.Name = getEnv("DB_NAME", "name")

	// Server configuration
	cfg.Server.Port = getEnv("PORT", "8080") // Using PORT instead of SERVER_PORT

	// JWT configuration
	cfg.JWT.AccessSecret = getEnv("ACCESS_TOKEN_SECRET", "default-secret-key")
	cfg.JWT.RefreshSecret = getEnv("REFRESH_TOKEN_SECRET", "default-secret-key")

	cfg.Log.Environment = getEnv("ENVIRONMENT", "Dev")
	cfg.Log.Version = getEnv("VERSION", "1.0.0")
	cfg.Log.ServiceName = getEnv("SERVICE_NAME", "auth-service")

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

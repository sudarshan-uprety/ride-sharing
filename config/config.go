package config

import (
	"log"
	"os"
	"strconv"

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
		Secret string
		Expiry int // in hours
	}
	Log struct {
		Environement string
		Version      string
		ServiceName  string
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
	cfg.JWT.Secret = getEnv("JWT_SECRET", "default-secret-key")
	expiry, err := strconv.Atoi(getEnv("JWT_EXPIRY", "24"))
	if err != nil {
		return nil, err
	}
	cfg.JWT.Expiry = expiry

	cfg.Log.Environement = getEnv("ENVIRONMENT", "Dev")
	cfg.Log.Environement = getEnv("VERSION", "1.0.0")
	cfg.Log.Environement = getEnv("SERVICE_NAME", "auth-service")

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

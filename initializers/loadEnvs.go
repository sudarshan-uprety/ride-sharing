package initializers

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	PORT   string
	DB_URL string
	SECRET string
)

func LoadEnvs() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	PORT = os.Getenv("PORT")
	DB_URL = os.Getenv("DB_URL")
	SECRET = os.Getenv("SECRET")
}

// load PostgreSQL, JWT, Appwrite, and AI API configurations
package config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

type Config struct {
	DatabaseURL  string
	JwtSecret    string
	GeminiAPIKey string
}

func LoadConfig() Config {
	return Config{
		DatabaseURL:  os.Getenv("DATABASE_URL"),
		JwtSecret:    os.Getenv("JWT_SECRET"),
		GeminiAPIKey: os.Getenv("GEMINI_API_KEY"),
	}
}

func ConnectDatabase(cfg Config) {
	database, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database: ", err)
	}

	DB = database
	fmt.Println("Database connected successfully")
}

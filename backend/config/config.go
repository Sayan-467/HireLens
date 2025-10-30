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
	DbUser       string
	DbPassword   string
	DbName       string
	DbPort       string
	DbHost       string
	JwtSecret    string
	GeminiAPIKey string
}

func LoadConfig() Config {
	return Config{
		DbUser:       os.Getenv("DB_USER"),
		DbPassword:   os.Getenv("DB_PASSWORD"),
		DbName:       os.Getenv("DB_NAME"),
		DbHost:       os.Getenv("DB_HOST"),
		DbPort:       os.Getenv("DB_PORT"),
		JwtSecret:    os.Getenv("JWT_SECRET"),
		GeminiAPIKey: os.Getenv("GEMINI_API_KEY"),
	}
}

func ConnectDatabase(cfg Config) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", cfg.DbHost, cfg.DbUser, cfg.DbPassword, cfg.DbName, cfg.DbPort)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database: ", err)
	}

	DB = database
	fmt.Println("Database connected successfully")
}

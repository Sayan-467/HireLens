package main

import (
	"backend/config"
	"backend/models"
	"backend/routes"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file only in local development (optional in production)
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables from system")
	}

	cfg := config.LoadConfig()
	config.ConnectDatabase(cfg)

	fmt.Println("JWT_SECRET from env:", os.Getenv("JWT_SECRET"))

	// auto migrate models
	err = config.DB.AutoMigrate(&models.User{}, &models.Resume{}, &models.JobRecommendation{})
	if err != nil {
		log.Fatal("Model migration failed", err)
	}
	log.Println("Database tables migrated successfully")

	router := gin.Default()

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	routes.SetupRoutes(router)

	log.Println("✅ Server starting on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v\n", err)
	}
}

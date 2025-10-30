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
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error handling .env", err)
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
	routes.SetupRoutes(router)
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v\n", err)
	}
}

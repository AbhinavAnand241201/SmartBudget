package main

import (
	"log"

	"smartbudget/api"
	"smartbudget/config"
	"smartbudget/db"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Set Gin mode based on environment
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize database connection
	if err := db.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.CloseDB()

	// Initialize router
	r := gin.Default()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"message": "SmartBudget API is running",
		})
	})

	// Initialize and register user routes
	userHandler := api.NewUserHandler()
	userHandler.RegisterRoutes(r)

	// Initialize and register insight routes
	insightHandler := api.NewInsightHandler()
	insightHandler.RegisterRoutes(r)

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"backend/database"
	"backend/handlers"
)

func main() {
	// Set release mode in production
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize database
	if err := database.Initialize(); err != nil {
		if err == database.ErrNoDatabasePath || err == database.ErrDatabaseNotFound {
			// Log the error but continue - we'll handle this in the API
			log.Printf("Database issue: %v", err)
		} else {
			log.Fatalf("Failed to initialize database: %v", err)
		}
	}
	defer database.Close()

	// Initialize Gin router
	router := gin.Default()

	// Configure CORS - allow all origins in production
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Content-Length"}
	corsConfig.ExposeHeaders = []string{"Content-Length"}
	corsConfig.AllowCredentials = true
	corsConfig.MaxAge = 12 * time.Hour
	router.Use(cors.New(corsConfig))

	// Set maximum multipart form size to 10MB
	router.MaxMultipartMemory = 10 << 20

	// IMPORTANT: API routes must be registered BEFORE setting up frontend serving
	// API routes
	api := router.Group("/api")
	{
		// History routes
		api.GET("/history", handlers.GetHistoryEntries)
		api.DELETE("/history", handlers.BatchDeleteHistoryEntries)

		// Config routes
		api.GET("/config/db-status", handlers.CheckDatabaseStatus)
	}

	// Setup frontend serving AFTER registering API routes
	ServeFrontend(router)

	// Determine port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	log.Printf("Server started at :%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

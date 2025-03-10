// main.go
package main

import (
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/dapoadedire/countries-states-cities-api/controller"
	"github.com/dapoadedire/countries-states-cities-api/database"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	database.InitDB()
	defer database.DB.Close()

	router := gin.Default()

	// Route for handling requests
	router.GET("/", controller.HandleWelcome)
	router.GET("/sync-data", controller.HandleSyncData)

	// Ensure the /data directory exists
	dataDir := "data"
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		if err := os.MkdirAll(dataDir, 0755); err != nil {
			log.Fatalf("Error creating directory: %v", err)
		}
	}

	// Get port from .env
	portStr := os.Getenv("PORT")
	if portStr == "" {
		portStr = "8080" // Default port if not set
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Invalid port number: %v", err)
	}

	// Start the server
	if err := router.Run(":" + strconv.Itoa(port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

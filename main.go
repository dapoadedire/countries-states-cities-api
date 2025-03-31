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
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	database.InitDB()
	defer database.DB.Close()

	router := gin.Default()

	syncEndpoint := os.Getenv("SYNC_ENDPOINT")
	if syncEndpoint == "" {
		syncEndpoint = "/sync-data"
	}

	router.GET("/", controller.HandleWelcome)
	router.POST(syncEndpoint, controller.HandleSyncAndPopulateData)
	router.GET("/countries", controller.HandleGetCountries)
	router.GET("/countries/:countryID", controller.HandleGetCountryByID)

	portStr := os.Getenv("PORT")
	if portStr == "" {
		portStr = "8080"
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Invalid port number: %v", err)
	}

	if err := router.Run(":" + strconv.Itoa(port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

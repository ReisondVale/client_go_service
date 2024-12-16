package main

import (
	"log"
	"project/internal/api"
	"project/internal/db"

	"github.com/gin-gonic/gin"
)

func main() {
	// Connection database
	database, err := db.Connect()
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	defer database.Close()

	// Define routes
	router := gin.Default()
	api.RegisterRoutes(router, database)

	// Start server
	log.Println("Server started on port 8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

package main

import (
	"log"
	"project/internal/api"
	"project/internal/db"
	"project/internal/repositories"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	requestCount uint64
	serverStart  time.Time
)

func main() {
	// Connection database
	database, err := db.Connect()
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	defer database.Close()

	serverStart = time.Now()
	// Define routes
	router := gin.Default()

	router.Use(func(c *gin.Context) {
		atomic.AddUint64(&requestCount, 1) // Increment the request counter
		c.Next()
	})

	clientRepo := &repositories.ClientRepository{DB: database}

	api.RegisterRoutes(router, clientRepo, &requestCount, &serverStart)

	// Start server
	log.Println("Server started on port 8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

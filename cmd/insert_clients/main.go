package main

import (
	"fmt"
	"log"
	"project/internal/db"
	"project/internal/repositories"
	"project/load_csv"
)

func main() {
	// Connect to the database
	dbConn, err := db.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbConn.Close()

	// Create a client repository
	clientRepo := &repositories.ClientRepository{DB: dbConn}

	// Define the CSV file path
	csvFilePath := "data/clients.csv"

	// Load clients from the CSV
	err = load_csv.LoadClientsFromCSV(csvFilePath, clientRepo)
	if err != nil {
		log.Fatalf("Failed to load clients from CSV: %v", err)
	}

	fmt.Println("Clients loaded successfully from CSV")
}

package load_csv

import (
	"encoding/csv"
	"fmt"
	"os"
	"project/internal/models"
	"project/internal/repositories"

	"github.com/klassmann/cpfcnpj"
)

func LoadClientsFromCSV(filePath string, repo *repositories.ClientRepository) error {
	// Open CSV
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open CSV file: %v", err)
	}
	defer file.Close()

	// Criar um leitor CSV
	reader := csv.NewReader(file)

	// Read all lines from CSV
	lines, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read CSV file: %v", err)
	}

	// Validate and process the lines
	for i, line := range lines {
		// Ignore header line
		if i == 0 {
			continue
		}

		// Check if each row has at least 2 columns
		if len(line) < 2 {
			return fmt.Errorf("invalid data at line %d: expected 2 columns, got %d", i+1, len(line))
		}

		// Create a new client from CSV data
		client := &models.Client{
			Name:     cpfcnpj.Clean(line[1]),
			CPF_CNPJ: cpfcnpj.Clean(line[0]),
		}

		// Insert client
		err = repo.Insert(client)
		if err != nil {
			return fmt.Errorf("failed to insert client at line %d: %v", i+1, err)
		}
	}

	return nil
}

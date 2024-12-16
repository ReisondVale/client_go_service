package repositories

import (
	"errors"
	"fmt"
	"project/internal/models"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/klassmann/cpfcnpj"
)

type ClientRepository struct {
	DB *sqlx.DB
}

func (r *ClientRepository) GetAll() ([]models.Client, error) {
	var clients []models.Client
	err := r.DB.Select(&clients, "SELECT * FROM clients ORDER BY name")
	return clients, err
}

func (r *ClientRepository) Insert(client *models.Client) error {
	// Check if CPF/CNPJ is valid and already exists in the database
	exists, err := r.Exists(client.CPF_CNPJ)
	if err != nil {
		return fmt.Errorf("failed to check CPF/CNPJ existence: %v", err)
	}
	if exists {
		return errors.New("CPF/CNPJ already registered")
	}

	// Start the transaction
	tx := r.DB.MustBegin()

	// Insert the new client into the database
	_, err = tx.NamedExec(`INSERT INTO clients (name, cpf_cnpj, blocklist) 
        VALUES (:name, :cpf_cnpj, :blocklist)`, client)
	if err != nil {
		if strings.Contains(err.Error(), "unique violation") {
			return errors.New("CPF/CNPJ already registered")
		}
		return fmt.Errorf("failed to insert client: %v", err)
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

func (r *ClientRepository) Exists(cpfCnpj string) (bool, error) {
	if len(cpfCnpj) == 11 {
		if !cpfcnpj.ValidateCPF(cpfCnpj) {
			return false, errors.New("invalid CPF")
		}
	} else if len(cpfCnpj) == 14 {
		if !cpfcnpj.ValidateCNPJ(cpfCnpj) {
			return false, errors.New("invalid CNPJ")
		}
	} else {
		return false, errors.New("invalid CPF/CNPJ length")
	}

	var count int
	err := r.DB.Get(&count, "SELECT COUNT(*) FROM clients WHERE cpf_cnpj = $1", cpfCnpj)
	if err != nil {
		return false, fmt.Errorf("failed to check if CPF/CNPJ exists: %v", err)
	}

	return count > 0, nil
}

func (r *ClientRepository) GetByName(name string) ([]models.Client, error) {
	var clients []models.Client

	// Query to search for clients by name (using ILIKE for case-insensitive search)
	err := r.DB.Select(&clients, "SELECT * FROM clients WHERE name ILIKE $1 ORDER BY name", "%"+name+"%")
	if err != nil {
		return nil, fmt.Errorf("failed to get clients by name: %v", err)
	}

	return clients, nil
}

package repositories

import (
	"errors"
	"testing"

	"project/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockClientRepository embute mock.Mock para simular métodos
type MockClientRepository struct {
	mock.Mock
}

func (m *MockClientRepository) GetAll() ([]models.Client, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Client), args.Error(1)
}

func (m *MockClientRepository) Insert(client *models.Client) error {
	args := m.Called(client)
	return args.Error(0)
}

func (m *MockClientRepository) Exists(cpfCnpj string) (bool, error) {
	args := m.Called(cpfCnpj)
	return args.Bool(0), args.Error(1)
}

func (m *MockClientRepository) GetByName(name string) ([]models.Client, error) {
	args := m.Called(name)
	return args.Get(0).([]models.Client), args.Error(1)
}

func TestGetAll_Success(t *testing.T) {
	mockRepo := &MockClientRepository{}
	mockClients := []models.Client{
		{ID: 1, Name: "John Doe", CPF_CNPJ: "12345678901", Blocklist: false},
		{ID: 2, Name: "Jane Doe", CPF_CNPJ: "98765432100", Blocklist: true},
	}

	mockRepo.On("GetAll").Return(mockClients, nil)

	clients, err := mockRepo.GetAll()
	assert.NoError(t, err)
	assert.Equal(t, len(mockClients), len(clients))
	mockRepo.AssertExpectations(t)
}

func TestInsert_ClientExists(t *testing.T) {
	mockRepo := &MockClientRepository{}
	client := &models.Client{ID: 1, Name: "John Doe", CPF_CNPJ: "12345678901"}

	mockRepo.On("Insert", client).Return(errors.New("CPF/CNPJ already registered"))

	err := mockRepo.Insert(client)
	assert.Error(t, err)
	assert.Equal(t, "CPF/CNPJ already registered", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestExists_InvalidCPF(t *testing.T) {
	mockRepo := &MockClientRepository{}
	mockRepo.On("Exists", "invalidCPF").Return(false, errors.New("invalid CPF"))

	exists, err := mockRepo.Exists("invalidCPF")
	assert.Error(t, err)
	assert.False(t, exists)
	mockRepo.AssertExpectations(t)
}

func TestGetByName_NoResults(t *testing.T) {
	mockRepo := &MockClientRepository{}
	mockRepo.On("GetByName", "Nonexistent").Return([]models.Client{}, nil)

	clients, err := mockRepo.GetByName("Nonexistent")
	assert.NoError(t, err)
	assert.Empty(t, clients)
	mockRepo.AssertExpectations(t)
}

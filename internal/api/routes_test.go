package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"project/internal/models"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock para o repositório de clientes
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

// Helper para configurar o roteador com mock
func setupRouter(mockRepo *MockClientRepository) *gin.Engine {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	requestCount := uint64(0)
	serverStart := time.Now()

	RegisterRoutes(router, mockRepo, &requestCount, &serverStart)
	return router
}

// Teste para rota /clients (GET)
func TestGetAllClients(t *testing.T) {
	mockRepo := new(MockClientRepository)
	mockRepo.On("GetAll").Return([]models.Client{
		{ID: 1, Name: "John Doe", CPF_CNPJ: "12345678901", Blocklist: false},
	}, nil)

	router := setupRouter(mockRepo)

	req, _ := http.NewRequest(http.MethodGet, "/clients", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var clients []models.Client
	err := json.Unmarshal(resp.Body.Bytes(), &clients)
	assert.NoError(t, err)
	assert.Len(t, clients, 1)
	mockRepo.AssertExpectations(t)
}

func TestGetAllClients_Error(t *testing.T) {
	mockRepo := new(MockClientRepository)
	mockRepo.On("GetAll").Return(nil, errors.New("database error"))

	router := setupRouter(mockRepo)

	req, _ := http.NewRequest(http.MethodGet, "/clients", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	mockRepo.AssertExpectations(t)
}

// Teste para rota /clients/exists/:cpfCnpj (GET)
func TestClientExists_Success(t *testing.T) {
	mockRepo := new(MockClientRepository)
	mockRepo.On("Exists", "12345678901").Return(true, nil)

	router := setupRouter(mockRepo)

	req, _ := http.NewRequest(http.MethodGet, "/clients/exists/12345678901", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	mockRepo.AssertExpectations(t)
}

func TestClientExists_InvalidCPF(t *testing.T) {
	mockRepo := new(MockClientRepository)
	mockRepo.On("Exists", "invalidCPF").Return(false, errors.New("invalid CPF"))

	router := setupRouter(mockRepo)

	req, _ := http.NewRequest(http.MethodGet, "/clients/exists/invalidCPF", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	mockRepo.AssertExpectations(t)
}

// Teste para rota /clients/search (GET)
func TestSearchClients_Success(t *testing.T) {
	mockRepo := new(MockClientRepository)
	mockRepo.On("GetByName", "John").Return([]models.Client{
		{ID: 1, Name: "John Doe", CPF_CNPJ: "12345678901", Blocklist: false},
	}, nil)

	router := setupRouter(mockRepo)

	req, _ := http.NewRequest(http.MethodGet, "/clients/search?name=John", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	mockRepo.AssertExpectations(t)
}

func TestSearchClients_NoName(t *testing.T) {
	mockRepo := new(MockClientRepository)
	router := setupRouter(mockRepo)

	req, _ := http.NewRequest(http.MethodGet, "/clients/search", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

// Teste para rota /clients (POST)
func TestInsertClient_Success(t *testing.T) {
	mockRepo := new(MockClientRepository)
	mockRepo.On("Insert", &models.Client{Name: "John Doe", CPF_CNPJ: "12345678901", Blocklist: false}).Return(nil)

	router := setupRouter(mockRepo)

	client := models.Client{Name: "John Doe", CPF_CNPJ: "12345678901", Blocklist: false}
	body, _ := json.Marshal(client)

	req, _ := http.NewRequest(http.MethodPost, "/clients", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)
	mockRepo.AssertExpectations(t)
}

func TestInsertClient_InvalidBody(t *testing.T) {
	mockRepo := new(MockClientRepository)
	router := setupRouter(mockRepo)

	req, _ := http.NewRequest(http.MethodPost, "/clients", strings.NewReader("invalid body"))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

// Teste para rota /status (GET)
func TestStatusEndpoint(t *testing.T) {
	mockRepo := new(MockClientRepository)
	router := setupRouter(mockRepo)

	req, _ := http.NewRequest(http.MethodGet, "/status", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

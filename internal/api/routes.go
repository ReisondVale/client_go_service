package api

import (
	"project/internal/models"
	"project/internal/repositories"
	"sync/atomic"
	"time"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func RegisterRoutes(router *gin.Engine, db *sqlx.DB, requestCount *uint64, serverStart *time.Time) {
	repo := &repositories.ClientRepository{DB: db}

	// Route to fetch all records
	router.GET("/clients", func(c *gin.Context) {
		clients, err := repo.GetAll()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, clients)
	})

	// Route to check if a CPF/CNPJ exists
	router.GET("/clients/exists/:cpfCnpj", func(c *gin.Context) {
		cpfCnpj := c.Param("cpfCnpj")
		exists, err := repo.Exists(cpfCnpj)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"exists": exists})
	})

	// Route to search for clients by name
	router.GET("/clients/search", func(c *gin.Context) {
		name := c.Query("name")
		if name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Name query parameter is required"})
			return
		}

		clients, err := repo.GetByName(name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, clients)
	})

	// Route to insert a new client
	router.POST("/clients", func(c *gin.Context) {
		var client models.Client
		if err := c.ShouldBindJSON(&client); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		err := repo.Insert(&client)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Client created successfully"})
	})

	// Route to verify application status
	router.GET("/status", func(c *gin.Context) {
		uptime := time.Since(*serverStart).Seconds()
		c.JSON(http.StatusOK, gin.H{
			"uptime_seconds": uptime,
			"request_count":  atomic.LoadUint64(requestCount),
		})
	})
}

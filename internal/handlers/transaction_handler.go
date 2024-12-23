package handlers

import (
	"net/http"
	"restaurant-stock/internal/models"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TransactionHandler struct {
	db *gorm.DB
}

func NewTransactionHandler(db *gorm.DB) *TransactionHandler {
	return &TransactionHandler{db: db}
}

type CreateTransactionRequest struct {
	TransType   string  `json:"transType"`
	ProductName string  `json:"productName"`
	Quantity    float64 `json:"quantity"`
	Unit        string  `json:"unit"`
}

func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	var req CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transaction := models.Transaction{
		TransType:     models.TransactionType(req.TransType),
		ProductName:   req.ProductName,
		Quantity:      req.Quantity,
		Unit:          req.Unit,
		TransactionAt: time.Now(),
	}

	if err := h.db.Create(&transaction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "İşlem kaydedilemedi"})
		return
	}

	c.JSON(http.StatusCreated, transaction)
}

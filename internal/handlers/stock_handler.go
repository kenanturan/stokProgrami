package handlers

import (
	"fmt"
	"net/http"
	"restaurant-stock/internal/models"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type StockHandler struct {
	db *gorm.DB
}

func NewStockHandler(db *gorm.DB) *StockHandler {
	return &StockHandler{db: db}
}

func (h *StockHandler) GetStockStatus(c *gin.Context) {
	var transactions []models.Transaction
	if err := h.db.Find(&transactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("Bulunan işlem sayısı: %d\n", len(transactions))

	// Stok durumunu hesapla
	stocks := make(map[string]float64)
	for _, t := range transactions {
		key := fmt.Sprintf("%s_%s", t.ProductName, t.Unit)

		// Giriş işlemi ise ekle, çıkış işlemi ise çıkar
		if t.TransType == models.TransactionTypeIn {
			stocks[key] += t.Quantity
		} else if t.TransType == models.TransactionTypeOut {
			stocks[key] -= t.Quantity
		}
	}

	// Sonuçları hazırla
	var result []gin.H
	for key, quantity := range stocks {
		parts := strings.Split(key, "_")
		productName := parts[0]
		unit := parts[1]

		result = append(result, gin.H{
			"productName": productName,
			"quantity":    quantity,
			"unit":        unit,
		})
	}

	c.JSON(http.StatusOK, result)
}

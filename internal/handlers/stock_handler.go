package handlers

import (
	"net/http"
	"restaurant-stock/internal/models"

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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Stok durumu alınamadı"})
		return
	}

	stocks := models.CalculateStock(transactions)

	// Map'i slice'a çevir
	var stockList []models.StockStatus
	for _, stock := range stocks {
		stockList = append(stockList, stock)
	}

	c.JSON(http.StatusOK, stockList)
}

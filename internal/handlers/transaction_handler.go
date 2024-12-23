package handlers

import (
	"fmt"
	"net/http"
	"restaurant-stock/internal/models"
	"strings"

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
	var transaction models.Transaction
	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transaction.TransType = models.TransactionType(strings.ToLower(string(transaction.TransType)))
	transaction.Unit = strings.ToLower(transaction.Unit)

	fmt.Printf("\n=== Yeni Stok Hareketi ===\n")
	fmt.Printf("Tür: %s\n", transaction.TransType)
	fmt.Printf("Ürün: %s\n", transaction.ProductName)
	fmt.Printf("Miktar: %.2f %s\n", transaction.Quantity, transaction.Unit)

	tx := h.db.Begin()

	if transaction.TransType == models.TransactionType("out") {
		fmt.Printf("\nÇıkış işlemi tespit edildi, reçete kontrolü yapılıyor...\n")

		var recipe models.Recipe
		result := tx.Debug().
			Preload("Ingredients").
			Where("LOWER(name) = LOWER(?)", transaction.ProductName).
			First(&recipe)

		if result.Error != nil {
			fmt.Printf("Reçete bulunamadı veya hata: %v\n", result.Error)
		} else {
			fmt.Printf("\nReçete bulundu: %s (ID: %d)\n", recipe.Name, recipe.ID)
			fmt.Printf("Malzeme sayısı: %d\n", len(recipe.Ingredients))

			// Her malzeme için stok çıkışı oluştur
			for _, ing := range recipe.Ingredients {
				if ing.Name != "" && ing.Quantity > 0 {
					ingredientQty := ing.Quantity * transaction.Quantity
					ingredientTx := models.Transaction{
						TransType:   models.TransactionType("out"),
						ProductName: ing.Name,
						Quantity:    ingredientQty,
						Unit:        strings.ToLower(ing.Unit),
					}

					fmt.Printf("\nMalzeme çıkışı oluşturuluyor:\n")
					fmt.Printf("- Ürün: %s\n", ingredientTx.ProductName)
					fmt.Printf("- Miktar: %.2f %s\n", ingredientTx.Quantity, ingredientTx.Unit)

					if err := tx.Debug().Create(&ingredientTx).Error; err != nil {
						tx.Rollback()
						fmt.Printf("Malzeme çıkışı yapılamadı: %v\n", err)
						c.JSON(http.StatusInternalServerError, gin.H{"error": "Malzeme çıkışı yapılamadı"})
						return
					}

					fmt.Printf("✓ Malzeme stoktan düşüldü: %s (%.2f %s)\n",
						ingredientTx.ProductName, ingredientTx.Quantity, ingredientTx.Unit)
				}
			}
		}
	}

	// Ana işlemi kaydet
	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		fmt.Printf("Transaction kaydedilemedi: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "İşlem kaydedilemedi"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		fmt.Printf("İşlem onaylanamadı: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "İşlem onaylanamadı"})
		return
	}

	fmt.Printf("\n✓ İşlem başarıyla kaydedildi\n")
	c.JSON(http.StatusCreated, transaction)
}

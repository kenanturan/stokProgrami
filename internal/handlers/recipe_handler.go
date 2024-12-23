package handlers

import (
	"fmt"
	"net/http"
	"restaurant-stock/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RecipeHandler struct {
	db *gorm.DB
}

func NewRecipeHandler(db *gorm.DB) *RecipeHandler {
	return &RecipeHandler{db: db}
}

// Reçete oluşturma isteği için struct
type CreateRecipeRequest struct {
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Ingredients []RecipeIngredientReq `json:"ingredients"`
}

type RecipeIngredientReq struct {
	Name     string  `json:"name"`
	Quantity float64 `json:"quantity"`
	Unit     string  `json:"unit"`
}

// Reçete oluşturma
func (h *RecipeHandler) CreateRecipe(c *gin.Context) {
	var req CreateRecipeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Debug için request'i yazdır
	fmt.Printf("Gelen istek: %+v\n", req)

	recipe := models.Recipe{
		Name:        req.Name,
		Description: req.Description,
	}

	// Transaction başlat
	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Önce reçeteyi kaydet
	if err := tx.Debug().Create(&recipe).Error; err != nil {
		tx.Rollback()
		fmt.Printf("Reçete kaydedilirken hata: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Reçete kaydedilemedi"})
		return
	}

	// Malzemeleri kaydet
	if len(req.Ingredients) > 0 {
		ingredients := make([]models.RecipeIngredient, len(req.Ingredients))
		for i, ing := range req.Ingredients {
			ingredients[i] = models.RecipeIngredient{
				RecipeID: recipe.ID,
				Name:     ing.Name,
				Quantity: ing.Quantity,
				Unit:     ing.Unit,
			}
		}

		if err := tx.Debug().Create(&ingredients).Error; err != nil {
			tx.Rollback()
			fmt.Printf("Malzemeler kaydedilirken hata: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Malzemeler kaydedilemedi"})
			return
		}
	}

	// Reçeteyi malzemeleriyle birlikte tekrar yükle
	if err := tx.Preload("Ingredients").First(&recipe, recipe.ID).Error; err != nil {
		tx.Rollback()
		fmt.Printf("Reçete yüklenirken hata: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Reçete yüklenemedi"})
		return
	}

	// Transaction'ı commit et
	if err := tx.Commit().Error; err != nil {
		fmt.Printf("Transaction commit hatası: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "İşlem tamamlanamadı"})
		return
	}

	fmt.Printf("Reçete başarıyla kaydedildi: %+v\n", recipe)
	c.JSON(http.StatusCreated, recipe)
}

// Tüm reçeteleri listele
func (h *RecipeHandler) ListRecipes(c *gin.Context) {
	var recipes []models.Recipe

	// Önce recipes tablosunu kontrol et
	var recipeCount int64
	if err := h.db.Debug().Model(&models.Recipe{}).Count(&recipeCount).Error; err != nil {
		fmt.Printf("Reçete sayısı alınırken hata: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Reçete sayısı alınamadı"})
		return
	}
	fmt.Printf("Veritabanında %d reçete bulundu\n", recipeCount)

	// Sonra recipe_ingredients tablosunu kontrol et
	var ingredientCount int64
	if err := h.db.Debug().Model(&models.RecipeIngredient{}).Count(&ingredientCount).Error; err != nil {
		fmt.Printf("Malzeme sayısı alınırken hata: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Malzeme sayısı alınamadı"})
		return
	}
	fmt.Printf("Veritabanında %d malzeme bulundu\n", ingredientCount)

	// Ana sorguyu çalıştır
	result := h.db.Debug().
		Preload("Ingredients").
		Find(&recipes)

	if result.Error != nil {
		fmt.Printf("Veritabanı hatası: %v\n", result.Error)
		fmt.Printf("SQL: %v\n", result.Statement.SQL.String())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Reçeteler alınamadı: %v", result.Error),
		})
		return
	}

	// Detaylı sonuçları yazdır
	fmt.Printf("\n=== Bulunan Reçeteler ===\n")
	for i, recipe := range recipes {
		fmt.Printf("\nReçete %d:\n", i+1)
		fmt.Printf("  ID: %d\n", recipe.ID)
		fmt.Printf("  Ad: %s\n", recipe.Name)
		fmt.Printf("  Açıklama: %s\n", recipe.Description)
		fmt.Printf("  Malzemeler (%d adet):\n", len(recipe.Ingredients))
		for j, ing := range recipe.Ingredients {
			fmt.Printf("    %d. %s (%.2f %s)\n", j+1, ing.Name, ing.Quantity, ing.Unit)
		}
	}

	c.JSON(http.StatusOK, recipes)
}

// Reçete detayını getir
func (h *RecipeHandler) GetRecipe(c *gin.Context) {
	id := c.Param("id")
	var recipe models.Recipe

	if err := h.db.Preload("Ingredients").First(&recipe, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Reçete bulunamadı"})
		return
	}

	c.JSON(http.StatusOK, recipe)
}

// Reçeteye göre stok kontrolü
func (h *RecipeHandler) CheckStock(c *gin.Context) {
	id := c.Param("id")
	var recipe models.Recipe

	if err := h.db.Preload("Ingredients").First(&recipe, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Reçete bulunamadı"})
		return
	}

	// Mevcut stok durumunu al
	var transactions []models.Transaction
	if err := h.db.Find(&transactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Stok durumu alınamadı"})
		return
	}

	stocks := models.CalculateStock(transactions)

	// Her malzeme için stok kontrolü
	var missingIngredients []map[string]interface{}
	for _, ingredient := range recipe.Ingredients {
		stock, exists := stocks[ingredient.Name]
		if !exists || stock.Quantity < ingredient.Quantity {
			var availableAmount float64
			if exists {
				availableAmount = stock.Quantity
			}

			missingIngredients = append(missingIngredients, map[string]interface{}{
				"name":            ingredient.Name,
				"required":        ingredient.Quantity,
				"available":       exists,
				"availableAmount": availableAmount,
				"unit":            ingredient.Unit,
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"recipe":             recipe,
		"missingIngredients": missingIngredients,
		"canCook":            len(missingIngredients) == 0,
	})
}

// Reçete silme
func (h *RecipeHandler) DeleteRecipe(c *gin.Context) {
	id := c.Param("id")
	fmt.Printf("\n=== DELETE İsteği Alındı ===\n")
	fmt.Printf("URL: %s\n", c.Request.URL.String())
	fmt.Printf("Method: %s\n", c.Request.Method)
	fmt.Printf("ID: %s\n", id)
	fmt.Printf("Headers: %v\n", c.Request.Header)

	// Debug için SQL sorgularını göster
	h.db = h.db.Debug()

	// Reçeteyi bul ve sil
	var recipe models.Recipe
	if err := h.db.First(&recipe, id).Error; err != nil {
		fmt.Printf("Reçete bulunamadı: %v\n", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Reçete bulunamadı"})
		return
	}

	fmt.Printf("Reçete bulundu: %+v\n", recipe)

	// Malzemeleri sil
	if err := h.db.Where("recipe_id = ?", id).Delete(&models.RecipeIngredient{}).Error; err != nil {
		fmt.Printf("Malzemeler silinemedi: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Malzemeler silinemedi"})
		return
	}

	// Reçeteyi sil
	if err := h.db.Delete(&recipe).Error; err != nil {
		fmt.Printf("Reçete silinemedi: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Reçete silinemedi"})
		return
	}

	fmt.Printf("Reçete başarıyla silindi: ID=%d\n", recipe.ID)
	c.JSON(http.StatusOK, gin.H{"message": "Reçete başarıyla silindi"})
}

// Reçeteyi üret ve stoktan düş
func (h *RecipeHandler) ProduceRecipe(c *gin.Context) {
	id := c.Param("id")
	fmt.Printf("\n=== Reçete Üretim İsteği ===\n")
	fmt.Printf("ID: %s\n", id)

	// Reçeteyi bul
	var recipe models.Recipe
	if err := h.db.Preload("Ingredients").First(&recipe, id).Error; err != nil {
		fmt.Printf("Reçete bulunamadı: %v\n", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Reçete bulunamadı"})
		return
	}

	// Transaction başlat
	tx := h.db.Begin()

	// Her malzeme için stoktan düşme işlemi yap
	for _, ingredient := range recipe.Ingredients {
		// Stok çıkış hareketi oluştur
		transaction := models.Transaction{
			TransType:   models.TransactionTypeOut,
			ProductName: ingredient.Name,
			Quantity:    ingredient.Quantity,
			Unit:        ingredient.Unit,
		}

		if err := tx.Create(&transaction).Error; err != nil {
			tx.Rollback()
			fmt.Printf("Stok hareketi oluşturulamadı: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Stok hareketi oluşturulamadı"})
			return
		}
	}

	// İşlemi onayla
	if err := tx.Commit().Error; err != nil {
		fmt.Printf("İşlem onaylanamadı: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "İşlem onaylanamadı"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Reçete üretildi ve stoktan düşüldü",
		"recipe":  recipe,
	})
}

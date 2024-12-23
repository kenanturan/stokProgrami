package main

import (
	"fmt"
	"log"
	"net/http"
	"restaurant-stock/internal/database"
	"restaurant-stock/internal/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := database.InitDB()
	if err != nil {
		log.Fatal("Veritabanı bağlantısı kurulamadı:", err)
	}

	// Debug modu aktifleştir
	gin.SetMode(gin.DebugMode)
	r := gin.Default()

	// Template dosyalarını yükle
	r.LoadHTMLGlob("templates/*")

	// Handler'ları oluştur
	reportHandler := handlers.NewReportHandler(db)
	transactionHandler := handlers.NewTransactionHandler(db)
	stockHandler := handlers.NewStockHandler(db)
	recipeHandler := handlers.NewRecipeHandler(db)

	// 404 handler'ı ekle
	r.NoRoute(func(c *gin.Context) {
		fmt.Printf("\n=== 404 Not Found ===\n")
		fmt.Printf("Method: %s\n", c.Request.Method)
		fmt.Printf("URL: %s\n", c.Request.URL.String())
		fmt.Printf("Headers: %v\n", c.Request.Header)
		c.JSON(http.StatusNotFound, gin.H{
			"error":  "Route not found",
			"path":   c.Request.URL.Path,
			"method": c.Request.Method,
		})
	})

	// Ana sayfa
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "layout.html", nil)
	})

	// API rotaları
	r.GET("/api/reports", reportHandler.GetTransactionReport)
	r.POST("/api/transactions", transactionHandler.CreateTransaction)
	r.GET("/api/stocks", stockHandler.GetStockStatus)

	// Reçete rotaları - sıralama önemli
	r.GET("/api/recipes", recipeHandler.ListRecipes)
	r.POST("/api/recipes", recipeHandler.CreateRecipe)
	r.GET("/api/recipes/:id/check-stock", recipeHandler.CheckStock)
	r.POST("/api/recipes/:id/produce", recipeHandler.ProduceRecipe)
	r.DELETE("/api/recipes/:id", recipeHandler.DeleteRecipe)
	r.GET("/api/recipes/:id", recipeHandler.GetRecipe)

	// Route'ları yazdır
	fmt.Println("\n=== Tanımlı Route'lar ===")
	for _, route := range r.Routes() {
		fmt.Printf("%s %s -> %T\n", route.Method, route.Path, route.HandlerFunc)
	}

	// Sunucuyu başlat
	fmt.Println("\nSunucu başlatılıyor...")
	r.Run(":8080")
}

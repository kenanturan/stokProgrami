package main

import (
	"log"
	"restaurant-stock/internal/database"
	"restaurant-stock/internal/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := database.InitDB()
	if err != nil {
		log.Fatal("Veritabanı bağlantısı kurulamadı:", err)
	}

	r := gin.Default()

	// Template dosyalarını yükle
	r.LoadHTMLGlob("templates/*")

	// Handler'ları oluştur
	reportHandler := handlers.NewReportHandler(db)
	transactionHandler := handlers.NewTransactionHandler(db)
	stockHandler := handlers.NewStockHandler(db)

	// Rotaları tanımla
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "layout.html", nil)
	})

	r.GET("/api/reports", reportHandler.GetTransactionReport)
	r.POST("/api/transactions", transactionHandler.CreateTransaction)
	r.GET("/api/stocks", stockHandler.GetStockStatus)

	// Sunucuyu başlat
	r.Run(":8080")
}

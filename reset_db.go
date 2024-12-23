package main

import (
	"fmt"
	"log"
	"restaurant-stock/internal/database"
	"restaurant-stock/internal/models"
)

func main() {
	db, err := database.InitDB()
	if err != nil {
		log.Fatal("Veritabanı bağlantısı kurulamadı:", err)
	}

	// Tabloları sil
	db.Migrator().DropTable(&models.Transaction{})
	db.Migrator().DropTable(&models.RecipeIngredient{})
	db.Migrator().DropTable(&models.Recipe{})

	// Tabloları yeniden oluştur
	db.AutoMigrate(&models.Transaction{})
	db.AutoMigrate(&models.Recipe{})
	db.AutoMigrate(&models.RecipeIngredient{})

	fmt.Println("Veritabanı başarıyla sıfırlandı!")
}

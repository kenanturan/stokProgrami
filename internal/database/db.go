package database

import (
	"restaurant-stock/internal/models"

	"fmt"

	"gorm.io/driver/postgres" // veya mysql için mysql driver
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	dsn := "host=localhost user=kenanturan dbname=restaurant_stock port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Debug modunu aktif et
	db = db.Debug()

	// Tabloları sil
	fmt.Println("Tabloları siliniyor...")
	db.Migrator().DropTable(&models.RecipeIngredient{})
	db.Migrator().DropTable(&models.Recipe{})
	db.Migrator().DropTable(&models.Transaction{})

	// Auto Migration
	fmt.Println("Tablolar oluşturuluyor...")

	fmt.Println("Recipe tablosu oluşturuluyor...")
	if err := db.AutoMigrate(&models.Recipe{}); err != nil {
		return nil, fmt.Errorf("Recipe migration error: %v", err)
	}

	fmt.Println("RecipeIngredient tablosu oluşturuluyor...")
	if err := db.AutoMigrate(&models.RecipeIngredient{}); err != nil {
		return nil, fmt.Errorf("RecipeIngredient migration error: %v", err)
	}

	fmt.Println("Transaction tablosu oluşturuluyor...")
	if err := db.AutoMigrate(&models.Transaction{}); err != nil {
		return nil, fmt.Errorf("Transaction migration error: %v", err)
	}

	fmt.Println("Veritabanı başarıyla hazırlandı!")
	return db, nil
}

package database

import (
	"restaurant-stock/internal/models"

	"gorm.io/driver/postgres" // veya mysql için mysql driver
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	// Sistem kullanıcı adınızı kullanın (örnek: kenanturan)
	dsn := "host=localhost user=kenanturan dbname=restaurant_stock port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto Migration
	err = db.AutoMigrate(&models.Transaction{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

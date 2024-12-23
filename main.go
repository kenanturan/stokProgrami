package main

import (
	"fmt"

	"github.com/your-project/database"
	"github.com/your-project/models"
	"gorm.io/gorm"
)

func resetDatabase(db *gorm.DB) {
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

func main() {
	// ... mevcut kodlar ...

	db := database.InitDB()

	// Veritabanını sıfırla
	resetDatabase(db)

	// ... mevcut kodlar ...
}

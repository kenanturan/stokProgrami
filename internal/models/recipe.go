package models

type Recipe struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"column:name"`
	Description string `gorm:"column:description"`
	// İlişkiyi daha açık tanımlayalım
	Ingredients []RecipeIngredient `gorm:"foreignKey:RecipeID;references:ID"`
}

type RecipeIngredient struct {
	ID       uint    `gorm:"primaryKey"`
	RecipeID uint    `gorm:"column:recipe_id;not null"`
	Name     string  `gorm:"column:name"`
	Quantity float64 `gorm:"column:quantity"`
	Unit     string  `gorm:"column:unit"`
}

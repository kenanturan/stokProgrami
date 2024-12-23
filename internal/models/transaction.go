package models

import (
	"time"
)

type TransactionType string

const (
	TransactionTypeIn  TransactionType = "IN"
	TransactionTypeOut TransactionType = "OUT"
)

type Transaction struct {
	ID            uint            `gorm:"primaryKey"`
	TransType     TransactionType `gorm:"column:trans_type"`
	ProductName   string          `gorm:"column:product_name"`
	Quantity      float64         `gorm:"column:quantity"`
	Unit          string          `gorm:"column:unit"`
	TransactionAt time.Time       `gorm:"column:transaction_at;default:CURRENT_TIMESTAMP"`
	CreatedAt     time.Time       `gorm:"column:created_at"`
	UpdatedAt     time.Time       `gorm:"column:updated_at"`
}

package models

import (
	"time"
)

type TransactionType string

const (
	TransactionTypeIn  TransactionType = "in"
	TransactionTypeOut TransactionType = "out"
)

type Transaction struct {
	ID            uint64          `json:"id" gorm:"primaryKey"`
	TransType     TransactionType `json:"transType"`
	ProductName   string          `json:"productName"`
	Quantity      float64         `json:"quantity"`
	Unit          string          `json:"unit"`
	TransactionAt time.Time       `json:"transactionAt" gorm:"default:CURRENT_TIMESTAMP"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

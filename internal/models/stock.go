package models

type StockStatus struct {
	ProductName string `gorm:"primaryKey"`
	Quantity    float64
	Unit        string
}

// Stok durumunu hesaplamak için yardımcı fonksiyon
func CalculateStock(transactions []Transaction) map[string]StockStatus {
	stocks := make(map[string]StockStatus)

	for _, t := range transactions {
		stock, exists := stocks[t.ProductName]
		if !exists {
			stock = StockStatus{
				ProductName: t.ProductName,
				Unit:        t.Unit,
			}
		}

		if t.TransType == TransactionTypeIn {
			stock.Quantity += t.Quantity
		} else {
			stock.Quantity -= t.Quantity
		}

		stocks[t.ProductName] = stock
	}

	return stocks
}

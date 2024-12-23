package handlers

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"restaurant-stock/internal/models"
	"time"

	"github.com/jung-kurt/gofpdf"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	db *gorm.DB
}

func NewReportHandler(db *gorm.DB) *ReportHandler {
	return &ReportHandler{db: db}
}

func (h *ReportHandler) GetTransactionReport(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	// Debug için
	fmt.Printf("Start Date: %s, End Date: %s\n", startDate, endDate)

	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz başlangıç tarihi"})
		return
	}

	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz bitiş tarihi"})
		return
	}

	// Bitiş tarihini günün sonuna ayarla
	end = end.Add(24 * time.Hour).Add(-time.Second)

	var transactions []models.Transaction
	result := h.db.Where("transaction_at BETWEEN ? AND ?", start, end).
		Order("transaction_at desc").
		Find(&transactions)

	// Debug için
	fmt.Printf("Query result: %+v\n", transactions)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Rapor oluşturulurken hata oluştu"})
		return
	}

	// Format seçeneğine göre yanıt döndürme
	format := c.Query("format")
	switch format {
	case "csv":
		generateCSV(c, transactions)
	case "pdf":
		generatePDF(c, transactions)
	default:
		c.JSON(http.StatusOK, transactions)
	}
}

func generateCSV(c *gin.Context, transactions []models.Transaction) {
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment;filename=rapor.csv")

	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	// Başlıkları yaz
	headers := []string{"İşlem ID", "İşlem Türü", "Ürün Adı", "Miktar", "Birim", "İşlem Tarihi"}
	writer.Write(headers)

	// Verileri yaz
	for _, t := range transactions {
		row := []string{
			fmt.Sprint(t.ID),
			string(t.TransType),
			t.ProductName,
			fmt.Sprint(t.Quantity),
			t.Unit,
			t.TransactionAt.Format("2006-01-02 15:04"),
		}
		writer.Write(row)
	}
}

func generatePDF(c *gin.Context, transactions []models.Transaction) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)

	// Başlık
	pdf.Cell(190, 10, "Stok Hareket Raporu")
	pdf.Ln(20)

	// Tablo başlıkları
	pdf.SetFont("Arial", "B", 12)
	headers := []string{"Ürün", "İşlem Türü", "Miktar", "Birim", "Tarih"}
	for _, header := range headers {
		pdf.Cell(38, 10, header)
	}
	pdf.Ln(-1)

	// Veriler
	pdf.SetFont("Arial", "", 12)
	for _, t := range transactions {
		pdf.Cell(38, 10, t.ProductName)
		pdf.Cell(38, 10, string(t.TransType))
		pdf.Cell(38, 10, fmt.Sprintf("%.2f", t.Quantity))
		pdf.Cell(38, 10, t.Unit)
		pdf.Cell(38, 10, t.TransactionAt.Format("2006-01-02"))
		pdf.Ln(-1)
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "attachment;filename=rapor.pdf")
	pdf.Output(c.Writer)
}

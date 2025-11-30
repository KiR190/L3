package handler

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"time"

	"sales-tracker/internal/models"

	"github.com/gocarina/gocsv"
	"github.com/wb-go/wbf/ginext"
)

// ItemCSV представляет структуру для экспорта в CSV
type ItemCSV struct {
	ID          string  `csv:"id"`
	Type        string  `csv:"type"`
	Amount      float64 `csv:"amount"`
	Currency    string  `csv:"currency"`
	CategoryID  string  `csv:"category_id"`
	Description string  `csv:"description"`
	OccurredAt  string  `csv:"occurred_at"`
	CreatedAt   string  `csv:"created_at"`
}

// ExportCSV экспортирует данные в CSV формат
func (h *Handler) ExportCSV(c *ginext.Context) {
	// Получаем параметры фильтрации
	fromStr := c.Query("from")
	toStr := c.Query("to")
	typeFilter := c.Query("type")

	// Получаем все записи
	items, err := h.service.GetItems(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	// Фильтруем данные на основе параметров
	var filteredItems []models.Item
	for _, item := range items {
		// Фильтр по дате
		if fromStr != "" {
			from, err := time.Parse(time.RFC3339, fromStr)
			if err == nil && item.OccurredAt.Before(from) {
				continue
			}
		}
		if toStr != "" {
			to, err := time.Parse(time.RFC3339, toStr)
			if err == nil && item.OccurredAt.After(to) {
				continue
			}
		}

		// Фильтр по типу
		if typeFilter != "" && string(item.Type) != typeFilter {
			continue
		}

		filteredItems = append(filteredItems, item)
	}

	// Конвертируем в CSV структуру
	csvItems := make([]ItemCSV, 0, len(filteredItems))
	for _, item := range filteredItems {
		categoryID := ""
		if item.CategoryID != nil {
			categoryID = *item.CategoryID
		}
		description := ""
		if item.Description != nil {
			description = *item.Description
		}

		csvItems = append(csvItems, ItemCSV{
			ID:          item.ID,
			Type:        string(item.Type),
			Amount:      float64(item.Amount) / 100, // Конвертируем копейки в рубли
			Currency:    item.Currency,
			CategoryID:  categoryID,
			Description: description,
			OccurredAt:  item.OccurredAt.Format(time.RFC3339),
			CreatedAt:   item.CreatedAt.Format(time.RFC3339),
		})
	}

	// Устанавливаем точку с запятой как разделитель для Excel
	gocsv.SetCSVWriter(func(out io.Writer) *gocsv.SafeCSVWriter {
		writer := csv.NewWriter(out)
		writer.Comma = ';'
		return gocsv.NewSafeCSVWriter(writer)
	})

	// Генерируем CSV
	csvContent, err := gocsv.MarshalString(&csvItems)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": "Failed to generate CSV"})
		return
	}

	// Устанавливаем заголовки для скачивания файла
	filename := fmt.Sprintf("sales-tracker-%s.csv", time.Now().Format("2006-01-02"))
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Type", "text/csv; charset=utf-8")

	// Добавляем BOM для корректного отображения кириллицы в Excel
	_, _ = c.Writer.Write([]byte("\xEF\xBB\xBF"))
	c.String(http.StatusOK, csvContent)
}

package handler

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"time"

	"warehouse/internal/models"

	"github.com/gocarina/gocsv"
	"github.com/wb-go/wbf/ginext"
)

type ItemHistoryCSV struct {
	ID        string `csv:"id"`
	ItemID    string `csv:"item_id"`
	Action    string `csv:"action"`
	UserID    string `csv:"user_id"`
	Username  string `csv:"username"`
	Role      string `csv:"role"`
	CreatedAt string `csv:"created_at"`
	OldData   string `csv:"old_data"`
	NewData   string `csv:"new_data"`
}

func (h *Handler) ExportItemHistoryCSV(c *ginext.Context) {
	itemID := c.Param("id")

	actionFilter := c.Query("action") // INSERT/UPDATE/DELETE
	userFilter := c.Query("username") // email in our case
	fromStr := c.Query("from")        // RFC3339
	toStr := c.Query("to")            // RFC3339

	var from *time.Time
	if fromStr != "" {
		v, err := time.Parse(time.RFC3339, fromStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, ginext.H{"error": "invalid from"})
			return
		}
		from = &v
	}

	var to *time.Time
	if toStr != "" {
		v, err := time.Parse(time.RFC3339, toStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, ginext.H{"error": "invalid to"})
			return
		}
		to = &v
	}

	history, err := h.service.GetItemHistory(c.Request.Context(), itemID, 10000, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	var filtered []models.ItemHistory
	for _, hst := range history {
		if actionFilter != "" && hst.Action != actionFilter {
			continue
		}
		if userFilter != "" && (hst.Username == nil || *hst.Username != userFilter) {
			continue
		}
		if from != nil && hst.CreatedAt.Before(*from) {
			continue
		}
		if to != nil && hst.CreatedAt.After(*to) {
			continue
		}
		filtered = append(filtered, hst)
	}

	out := make([]ItemHistoryCSV, 0, len(filtered))
	for _, hst := range filtered {
		userID := ""
		if hst.UserID != nil {
			userID = *hst.UserID
		}
		username := ""
		if hst.Username != nil {
			username = *hst.Username
		}
		role := ""
		if hst.Role != nil {
			role = *hst.Role
		}

		out = append(out, ItemHistoryCSV{
			ID:        hst.ID,
			ItemID:    hst.ItemID,
			Action:    hst.Action,
			UserID:    userID,
			Username:  username,
			Role:      role,
			CreatedAt: hst.CreatedAt.Format(time.RFC3339),
			OldData:   string(hst.OldData),
			NewData:   string(hst.NewData),
		})
	}

	// Устанавливаем точку с запятой как разделитель для Excel
	gocsv.SetCSVWriter(func(out io.Writer) *gocsv.SafeCSVWriter {
		writer := csv.NewWriter(out)
		writer.Comma = ';'
		return gocsv.NewSafeCSVWriter(writer)
	})

	// Генерируем CSV
	csvContent, err := gocsv.MarshalString(&out)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": "Failed to generate CSV"})
		return
	}

	// Устанавливаем заголовки для скачивания файла
	filename := fmt.Sprintf("item-history-%s-%s.csv", itemID, time.Now().Format("2006-01-02"))
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Type", "text/csv; charset=utf-8")

	// Добавляем BOM для корректного отображения кириллицы в Excel
	_, _ = c.Writer.Write([]byte("\xEF\xBB\xBF"))
	c.String(http.StatusOK, csvContent)
}

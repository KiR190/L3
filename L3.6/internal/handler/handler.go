package handler

import (
	"net/http"

	"sales-tracker/internal/models"
	"sales-tracker/internal/service"

	"github.com/wb-go/wbf/ginext"
)

type Handler struct {
	service service.ItemService
}

func NewHandler(s service.ItemService) *Handler {
	return &Handler{service: s}
}

func (h *Handler) CreateItem(c *ginext.Context) {
	var body models.ItemCreate
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": err.Error()})
		return
	}

	// Валидация
	if err := ValidateStruct(&body); err != nil {
		RespondWithValidationError(c, err)
		return
	}

	item, err := h.service.CreateItem(c, body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, item)
}

func (h *Handler) GetItems(c *ginext.Context) {
	items, err := h.service.GetItems(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *Handler) UpdateItem(c *ginext.Context) {
	id := c.Param("id")

	var body models.ItemUpdate
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": err.Error()})
		return
	}

	// Валидация
	if err := ValidateStruct(&body); err != nil {
		RespondWithValidationError(c, err)
		return
	}

	item, err := h.service.UpdateItem(c, id, body)
	if err != nil {
		c.JSON(http.StatusNotFound, ginext.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *Handler) DeleteItem(c *ginext.Context) {
	id := c.Param("id")

	if err := h.service.DeleteItem(c, id); err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) GetAnalytics(c *ginext.Context) {
	from := c.Query("from")
	to := c.Query("to")

	result, err := h.service.GetAnalytics(c.Request.Context(), from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ginext.H{
		"sum":    float64(result.Sum) / 100,
		"avg":    result.Avg / 100,
		"count":  result.Count,
		"median": result.Median / 100,
		"p90":    result.P90 / 100,
	})
}

package handler

import (
	"net/http"

	"shortener/internal/service"

	"github.com/wb-go/wbf/ginext"
)

type Handler struct {
	service *service.ShortenerService
}

func NewURLHandler(s *service.ShortenerService) *Handler {
	return &Handler{service: s}
}

// POST /shorten
func (h *Handler) Shorten(c *ginext.Context) {
	var req struct {
		URL    string `json:"url"`
		Custom string `json:"custom,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "invalid request"})
		return
	}

	ctx := c.Request.Context()

	shortURL, err := h.service.Create(ctx, req.URL, req.Custom)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, shortURL)
}

// GET /s/:short_url
func (h *Handler) Redirect(c *ginext.Context) {
	shortCode := c.Param("short_url")
	ctx := c.Request.Context()

	originalURL, err := h.service.Resolve(ctx, shortCode, c.Request.UserAgent())
	if err != nil {
		c.JSON(http.StatusNotFound, ginext.H{"error": "short url not found"})
		return
	}

	c.Redirect(http.StatusFound, originalURL)
}

// GET /analytics/:short_url
func (h *Handler) Analytics(c *ginext.Context) {
	shortCode := c.Param("short_url")
	ctx := c.Request.Context()

	stats, err := h.service.GetAnalytics(ctx, shortCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ginext.H{"analytics": stats})
}

// GET /analytics/latest
func (h *Handler) Latest(c *ginext.Context) {
	ctx := c.Request.Context()

	items, err := h.service.ListLatest(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, items)
}

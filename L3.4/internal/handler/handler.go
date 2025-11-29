package handler

import (
	"net/http"
	"strconv"

	"image-processor/internal/models"
	"image-processor/internal/service"

	"github.com/wb-go/wbf/ginext"
)

type Handler struct {
	service service.ImageService
}

func NewHandler(s service.ImageService) *Handler {
	return &Handler{service: s}
}

// POST /upload
func (h *Handler) UploadImage(c *ginext.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "file is required"})
		return
	}

	// параметры задачи
	taskType := c.DefaultPostForm("type", "resize")

	width, _ := strconv.Atoi(c.DefaultPostForm("width", "800"))
	height, _ := strconv.Atoi(c.DefaultPostForm("height", "800"))
	watermarkText := c.DefaultPostForm("text", "Watermark")

	params := models.TaskParams{
		Width:     width,
		Height:    height,
		Watermark: watermarkText,
	}

	id, err := h.service.UploadAndEnqueue(c, file, taskType, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ginext.H{"id": id})
}

// GET /image/:id
func (h *Handler) GetImage(c *ginext.Context) {
	id := c.Param("id")

	data, err := h.service.GetImage(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, ginext.H{"error": err.Error()})
		return
	}

	c.Data(http.StatusOK, "image/jpeg", data)
}

// DELETE /image/:id
func (h *Handler) DeleteImage(c *ginext.Context) {
	id := c.Param("id")

	err := h.service.DeleteImage(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// GET /status/:id
func (h *Handler) GetStatus(c *ginext.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "id required"})
		return
	}

	task, err := h.service.GetStatus(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, ginext.H{"error": "status not found"})
		return
	}

	c.JSON(http.StatusOK, ginext.H{
		"status": task.Status,
		"result": "/image/" + task.ImageID,
	})
}

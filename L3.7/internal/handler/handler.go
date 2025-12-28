package handler

import (
	"net/http"
	"strconv"
	"strings"

	"warehouse/internal/auth"
	"warehouse/internal/models"
	"warehouse/internal/service"

	"github.com/wb-go/wbf/ginext"
)

type Handler struct {
	service *service.Service
}

func NewHandler(s *service.Service) *Handler {
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

	item, err := h.service.CreateItem(c.Request.Context(), body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, item)
}

func (h *Handler) GetItems(c *ginext.Context) {
	limit := 50
	if limitStr := c.Query("limit"); limitStr != "" {
		v, err := strconv.Atoi(limitStr)
		if err != nil || v <= 0 {
			c.JSON(http.StatusBadRequest, ginext.H{"error": "invalid limit"})
			return
		}
		limit = v
	}

	offset := 0
	if offsetStr := c.Query("offset"); offsetStr != "" {
		v, err := strconv.Atoi(offsetStr)
		if err != nil || v < 0 {
			c.JSON(http.StatusBadRequest, ginext.H{"error": "invalid offset"})
			return
		}
		offset = v
	}

	items, err := h.service.GetItems(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, items)
}

func (h *Handler) GetItem(c *ginext.Context) {
	id := c.Param("id")

	item, err := h.service.GetItem(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ginext.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *Handler) UpdateItem(c *ginext.Context) {
	id := c.Param("id")

	var body models.ItemUpdate
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": err.Error()})
		return
	}

	if err := ValidateStruct(&body); err != nil {
		RespondWithValidationError(c, err)
		return
	}

	item, err := h.service.UpdateItem(c.Request.Context(), id, body)
	if err != nil {
		c.JSON(http.StatusNotFound, ginext.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *Handler) DeleteItem(c *ginext.Context) {
	id := c.Param("id")

	if err := h.service.DeleteItem(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// Auth

func (h *Handler) RegisterUser(c *ginext.Context) {
	var req models.UserRegister
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": err.Error()})
		return
	}
	if err := ValidateStruct(&req); err != nil {
		RespondWithValidationError(c, err)
		return
	}

	resp, err := h.service.Register(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *Handler) LoginUser(c *ginext.Context) {
	var req models.UserLogin
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": err.Error()})
		return
	}
	if err := ValidateStruct(&req); err != nil {
		RespondWithValidationError(c, err)
		return
	}

	resp, err := h.service.Login(c.Request.Context(), req)
	if err != nil {
		// Only return 401 for bad credentials. Unexpected errors (DB/schema/etc) should be 500,
		// otherwise the frontend will treat them as auth failures.
		if strings.Contains(err.Error(), "invalid email or password") {
			c.JSON(http.StatusUnauthorized, ginext.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetCurrentUser(c *ginext.Context) {
	userID, ok := auth.GetUserID(c.Request.Context())
	if !ok {
		c.JSON(http.StatusUnauthorized, ginext.H{"error": "unauthorized"})
		return
	}

	user, err := h.service.GetUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// Admin

func (h *Handler) UpdateUserRole(c *ginext.Context) {
	userID := c.Param("id")
	var req models.UpdateUserRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": err.Error()})
		return
	}
	if err := ValidateStruct(&req); err != nil {
		RespondWithValidationError(c, err)
		return
	}

	if err := h.service.UpdateUserRole(c.Request.Context(), userID, models.UserRole(req.Role)); err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// History

func (h *Handler) GetItemHistory(c *ginext.Context) {
	itemID := c.Param("id")

	limit := 50
	if limitStr := c.Query("limit"); limitStr != "" {
		v, err := strconv.Atoi(limitStr)
		if err != nil || v <= 0 {
			c.JSON(http.StatusBadRequest, ginext.H{"error": "invalid limit"})
			return
		}
		limit = v
	}

	offset := 0
	if offsetStr := c.Query("offset"); offsetStr != "" {
		v, err := strconv.Atoi(offsetStr)
		if err != nil || v < 0 {
			c.JSON(http.StatusBadRequest, ginext.H{"error": "invalid offset"})
			return
		}
		offset = v
	}

	history, err := h.service.GetItemHistory(c.Request.Context(), itemID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, history)
}

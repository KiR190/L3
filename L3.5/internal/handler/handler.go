package handler

import (
	"net/http"
	"strings"

	"event-booker/internal/auth"
	"event-booker/internal/models"
	"event-booker/internal/repository"
	"event-booker/internal/service"

	"github.com/wb-go/wbf/ginext"
)

type Handler struct {
	service       *service.Service
	defaultTimeout int
}

func NewHandler(svc *service.Service, defaultTimeout int) *Handler {
	return &Handler{
		service:       svc,
		defaultTimeout: defaultTimeout,
	}
}

// CreateEvent handles POST /events
func (h *Handler) CreateEvent(c *ginext.Context) {
	var req models.EventCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": err.Error()})
		return
	}

	event, err := h.service.CreateEvent(c.Request.Context(), req, h.defaultTimeout)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, event)
}

// GetEvent handles GET /events/:id
func (h *Handler) GetEvent(c *ginext.Context) {
	id := c.Param("id")

	eventDetail, err := h.service.GetEvent(c.Request.Context(), id)
	if err != nil {
		if err == repository.ErrEventNotFound {
			c.JSON(http.StatusNotFound, ginext.H{"error": "event not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, eventDetail)
}

// GetEvents handles GET /events
func (h *Handler) GetEvents(c *ginext.Context) {
	events, err := h.service.GetEvents(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, events)
}

// UpdateEvent handles PUT /events/:id (admin only)
func (h *Handler) UpdateEvent(c *ginext.Context) {
	eventID := c.Param("id")

	var req models.EventUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "Invalid request: " + err.Error()})
		return
	}

	event, err := h.service.UpdateEvent(c.Request.Context(), eventID, req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, ginext.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, event)
}

// Auth handlers

// RegisterUser handles POST /auth/register
func (h *Handler) RegisterUser(c *ginext.Context) {
	var req models.UserRegister
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": err.Error()})
		return
	}

	authResp, err := h.service.Register(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, authResp)
}

// LoginUser handles POST /auth/login
func (h *Handler) LoginUser(c *ginext.Context) {
	var req models.UserLogin
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": err.Error()})
		return
	}

	authResp, err := h.service.Login(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ginext.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, authResp)
}

// GetCurrentUser handles GET /auth/me
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

// LinkTelegram handles POST /auth/telegram
func (h *Handler) LinkTelegram(c *ginext.Context) {
	userID, ok := auth.GetUserID(c.Request.Context())
	if !ok {
		c.JSON(http.StatusUnauthorized, ginext.H{"error": "unauthorized"})
		return
	}

	var req models.LinkTelegramRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "Invalid request: " + err.Error()})
		return
	}

	err := h.service.LinkTelegram(c.Request.Context(), userID, req.TelegramUsername)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ginext.H{
		"message": "Telegram username linked and set as preferred notification. Please send /register to the bot to complete setup.",
	})
}

// UpdateProfile handles PUT /auth/me
func (h *Handler) UpdateProfile(c *ginext.Context) {
	userID, ok := auth.GetUserID(c.Request.Context())
	if !ok {
		c.JSON(http.StatusUnauthorized, ginext.H{"error": "unauthorized"})
		return
	}

	var req models.UpdateProfile
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "Invalid request: " + err.Error()})
		return
	}

	user, err := h.service.UpdateProfile(c.Request.Context(), userID, req)
	if err != nil {
		if strings.Contains(err.Error(), "email already in use") {
			c.JSON(http.StatusConflict, ginext.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.UserResponse{
		ID:                   user.ID,
		Email:                user.Email,
		Role:                 string(user.Role),
		TelegramUsername:     user.TelegramUsername,
		TelegramRegistered:   user.TelegramChatID != nil,
		PreferredNotification: user.PreferredNotification,
	})
}

// CreateBooking handles POST /events/:id/book
func (h *Handler) CreateBooking(c *ginext.Context) {
	eventID := c.Param("id")

	// Get user ID from context (set by auth middleware)
	userID, ok := auth.GetUserID(c.Request.Context())
	if !ok {
		c.JSON(http.StatusUnauthorized, ginext.H{"error": "unauthorized"})
		return
	}

	var req models.BookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": err.Error()})
		return
	}

	booking, err := h.service.CreateBooking(c.Request.Context(), eventID, userID, req)
	if err != nil {
		if err == repository.ErrEventNotFound {
			c.JSON(http.StatusNotFound, ginext.H{"error": "event not found"})
			return
		}
		if err == repository.ErrInsufficientSeats {
			c.JSON(http.StatusConflict, ginext.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, booking)
}

// ConfirmBooking handles POST /bookings/:id/confirm
func (h *Handler) ConfirmBooking(c *ginext.Context) {
	bookingID := c.Param("id")

	// Get user ID from context (set by auth middleware)
	userID, ok := auth.GetUserID(c.Request.Context())
	if !ok {
		c.JSON(http.StatusUnauthorized, ginext.H{"error": "unauthorized"})
		return
	}

	booking, err := h.service.ConfirmBooking(c.Request.Context(), bookingID, userID)
	if err != nil {
		if err == repository.ErrBookingNotFound {
			c.JSON(http.StatusNotFound, ginext.H{"error": "booking not found"})
			return
		}
		if err == repository.ErrBookingAlreadyPaid {
			c.JSON(http.StatusConflict, ginext.H{"error": "booking already paid"})
			return
		}
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, booking)
}

// GetMyBookings handles GET /bookings
func (h *Handler) GetMyBookings(c *ginext.Context) {
	// Get user ID from context (set by auth middleware)
	userID, ok := auth.GetUserID(c.Request.Context())
	if !ok {
		c.JSON(http.StatusUnauthorized, ginext.H{"error": "unauthorized"})
		return
	}

	bookings, err := h.service.GetUserBookings(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, bookings)
}

// GetBooking handles GET /bookings/:id
func (h *Handler) GetBooking(c *ginext.Context) {
	bookingID := c.Param("id")

	booking, err := h.service.GetBooking(c.Request.Context(), bookingID)
	if err != nil {
		if err == repository.ErrBookingNotFound {
			c.JSON(http.StatusNotFound, ginext.H{"error": "booking not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, ginext.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, booking)
}


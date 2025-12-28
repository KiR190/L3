package handler

import (
	"net/http"
	"time"

	"event-booker/internal/auth"

	"github.com/wb-go/wbf/ginext"
)

// NewRouter creates and configures the Gin router
func NewRouter(h *Handler, authService *auth.AuthService, adminMiddleware *auth.AdminMiddleware) *ginext.Engine {
	router := ginext.New("")
	router.Use(ginext.Logger())
	router.Use(ginext.Recovery())

	// Health check
	router.GET("/health", func(c *ginext.Context) {
		c.JSON(http.StatusOK, ginext.H{"status": "ok"})
	})

	// Public auth routes
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/register", h.RegisterUser)
		authGroup.POST("/login", h.LoginUser)
	}

	// Protected auth routes (require authentication)
	authProtected := router.Group("/auth")
	authProtected.Use(auth.RequireAuth(authService))
	{
		authProtected.GET("/me", h.GetCurrentUser)
		authProtected.PUT("/me", h.UpdateProfile)
		authProtected.POST("/telegram", h.LinkTelegram)
	}

	// Public event routes (anyone can view events)
	router.GET("/events", h.GetEvents)
	router.GET("/events/:id", h.GetEvent)

	// Admin-only event routes (require admin role)
	adminRoutes := router.Group("")
	adminRoutes.Use(auth.RequireAuth(authService))
	adminRoutes.Use(adminMiddleware.RequireAdmin())
	{
		adminRoutes.POST("/events", h.CreateEvent)
		adminRoutes.PUT("/events/:id", h.UpdateEvent)
	}

	// Protected booking routes (require authentication)
	protectedRoutes := router.Group("")
	protectedRoutes.Use(auth.RequireAuth(authService))
	{
		// Create booking (requires auth)
		protectedRoutes.POST("/events/:id/book", h.CreateBooking)
		
		// Booking operations (require auth)
		protectedRoutes.GET("/bookings", h.GetMyBookings)
		protectedRoutes.GET("/bookings/:id", h.GetBooking)
		protectedRoutes.POST("/bookings/:id/confirm", h.ConfirmBooking)
	}

	return router
}

// NewHTTPServer creates an HTTP server with the given router
func NewHTTPServer(port string, router http.Handler) *http.Server {
	return &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}


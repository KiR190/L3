package handler

import (
	"net/http"
	"time"

	"warehouse/internal/auth"
	"warehouse/internal/models"

	"github.com/wb-go/wbf/ginext"
)

// NewRouter создает Gin-роутер с маршрутами
func NewRouter(h *Handler, authService *auth.AuthService) *ginext.Engine {
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

	// Protected auth routes
	authProtected := router.Group("/auth")
	authProtected.Use(auth.RequireAuth(authService))
	{
		authProtected.GET("/me", h.GetCurrentUser)
	}

	// Items (protected)
	items := router.Group("/items")
	items.Use(auth.RequireAuth(authService))
	{
		// Read-only for viewer/manager/admin
		items.Use(auth.RequireAnyRole(string(models.RoleViewer), string(models.RoleManager), string(models.RoleAdmin)))
		items.GET("/", h.GetItems)
		items.GET("/:id", h.GetItem)
		items.GET("/:id/history", h.GetItemHistory)
		items.GET("/:id/history/export.csv", h.ExportItemHistoryCSV)

		// Edit for manager/admin
		editor := items.Group("")
		editor.Use(auth.RequireAnyRole(string(models.RoleManager), string(models.RoleAdmin)))
		editor.POST("/", h.CreateItem)
		editor.PUT("/:id", h.UpdateItem)

		// Delete only for admin
		adminOnly := items.Group("")
		adminOnly.Use(auth.RequireAnyRole(string(models.RoleAdmin)))
		adminOnly.DELETE("/:id", h.DeleteItem)
	}

	// Admin endpoints
	adminRoutes := router.Group("/admin")
	adminRoutes.Use(auth.RequireAuth(authService))
	adminRoutes.Use(auth.RequireAnyRole(string(models.RoleAdmin)))
	{
		adminRoutes.PUT("/users/:id/role", h.UpdateUserRole)
	}

	return router
}

func NewHTTPServer(port string, router http.Handler) *http.Server {
	return &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}

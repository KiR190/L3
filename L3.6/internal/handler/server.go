package handler

import (
	"net/http"
	"time"

	"github.com/wb-go/wbf/ginext"
)

// NewRouter создает Gin-роутер с маршрутами
func NewRouter(h *Handler) *ginext.Engine {
	router := ginext.New("")
	router.Use(ginext.Logger())
	router.Use(ginext.Recovery())

	// Роуты
	api := router.Group("/items")
	{
		api.POST("/", h.CreateItem)
		api.GET("/", h.GetItems)
		api.PUT("/:id", h.UpdateItem)
		api.DELETE("/:id", h.DeleteItem)
	}

	router.GET("/analytics", h.GetAnalytics)
	router.GET("/export/csv", h.ExportCSV)

	router.GET("/", func(c *ginext.Context) {
		c.File("./internal/handler/static/index.html")
	})

	router.Engine.Static("/static", "./internal/handler/static")

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

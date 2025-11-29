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

	router.POST("/upload", h.UploadImage)

	// Роуты
	api := router.Group("/image")
	{
		api.GET("/:id", h.GetImage)
		api.DELETE("/:id", h.DeleteImage)
	}

	router.GET("/", func(c *ginext.Context) {
		c.File("./internal/handler/static/index.html")
	})

	router.GET("/status/:id", h.GetStatus)

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

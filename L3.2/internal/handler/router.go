package handler

import (
	"github.com/wb-go/wbf/ginext"
)

func SetupRouter(e *ginext.Engine, h *Handler) {
	api := e.Group("")

	// POST /shorten — создание новой короткой ссылки
	api.POST("/shorten", h.Shorten)

	// GET /s/:short_url — переход по короткой ссылке
	api.GET("/s/:short_url", h.Redirect)

	// GET /analytics/:short_url — получение аналитики
	api.GET("/analytics/:short_url", h.Analytics)
	api.GET("/analytics/latest", h.Latest)

	api.GET("/", func(c *ginext.Context) {
		c.File("./internal/handler/static/index.html")
	})
	api.Static("/static", "./internal/handler/static")
}

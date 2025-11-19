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
	api := router.Group("/comments")
	{
		api.POST("", h.CreateComment)
		api.GET("", h.GetComments)
		api.DELETE("/:id", h.DeleteComment)
		api.GET("/search", h.SearchComments)
	}

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

package server

import (
	"net/http"
	"time"

	"delayed-notifier/config"
	"delayed-notifier/internal/handler"

	"github.com/gin-gonic/gin"
)

// NewRouter создает Gin-роутер с маршрутами
func NewRouter(notifHandler *handler.NotificationHandler) *gin.Engine {
	router := gin.Default()

	// Роуты уведомлений
	api := router.Group("/notify")
	{
		api.POST("", notifHandler.CreateNotification)
		api.GET("", notifHandler.ListNotifications)
		api.GET("/:id", notifHandler.GetNotification)
		api.DELETE("/:id", notifHandler.CancelNotification)
	}

	router.GET("/", func(c *gin.Context) {
		c.File("./internal/handler/static/index.html")
	})

	router.Static("/static", "./internal/handler/static")

	return router
}

func NewHTTPServer(cfg *config.Config, router http.Handler) *http.Server {
	return &http.Server{
		Addr:         ":" + cfg.HTTPPort,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}

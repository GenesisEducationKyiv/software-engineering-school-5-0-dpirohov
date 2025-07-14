package server

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"weatherApi/internal/server/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	api := r.Group("/api/v1")
	{
		weatherHandler := routes.NewWeatherHandler(s.WeatherService)
		api.GET("/health", s.healthHandler)
		api.GET("/weather", weatherHandler.GetWeather)

		subscriptionHandler := routes.NewSubscriptionHandler(s.SubscriptionService)
		api.POST("/subscribe", subscriptionHandler.Subscribe)
		api.GET("/confirm/:token", subscriptionHandler.ConfirmSubscription)
		api.GET("/unsubscribe/:token", subscriptionHandler.Unsubscribe)
	}

	webDir := filepath.Join(s.config.RootDir, "web")

	r.Static("/assets", filepath.Join(webDir, "assets"))

	// serve static page for unknown routes
	r.NoRoute(func(c *gin.Context) {
		indexPath := filepath.Join(webDir, "index.html")

		if _, err := os.Stat(indexPath); os.IsNotExist(err) {
			c.String(http.StatusNotFound, "index.html not found")
			return
		}

		c.File(indexPath)
	})

	return r
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.HealthCheckService.Health())
}

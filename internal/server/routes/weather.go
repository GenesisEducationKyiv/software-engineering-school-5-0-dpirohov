package routes

import (
	"net/http"
	"weatherApi/internal/logger"

	"weatherApi/internal/service/weather"

	"github.com/gin-gonic/gin"
)

type WeatherHandler struct {
	service *weather.Service
}

func NewWeatherHandler(weatherService *weather.Service) *WeatherHandler {
	return &WeatherHandler{
		service: weatherService,
	}
}

func (h *WeatherHandler) GetWeather(c *gin.Context) {
	log := logger.FromContext(c.Request.Context())

	city := c.Query("city")
	log.Info().Msgf("Handling get weather for %s", city)
	if city == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "city is required"})
		return
	}

	response, err := h.service.GetWeather(c.Request.Context(), city)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to get weather for %s", city)
		c.AbortWithStatusJSON(err.Code, gin.H{"error": err.Message})
		return
	}

	c.JSON(http.StatusOK, response)
}

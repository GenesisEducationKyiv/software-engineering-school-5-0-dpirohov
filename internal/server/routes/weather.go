package routes

import (
	"weatherApi/internal/service/weather"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WeatherHandler struct {
	service *weather.WeatherService
}

func NewWeatherHandler(weatherService *weather.WeatherService) *WeatherHandler {
	return &WeatherHandler{
		service: weatherService,
	}
}

func (h *WeatherHandler) GetWeather(c *gin.Context) {
	city := c.Query("city")

	if city == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "city is required"})
		return
	}

	weather, err := h.service.GetWeather(city)

	if err != nil {
		c.AbortWithStatusJSON(err.Code, gin.H{"error": err.Message})
		return
	}

	c.JSON(http.StatusOK, weather)
}

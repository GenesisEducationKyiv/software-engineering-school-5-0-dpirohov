package routes

import (
	"net/http"

	"weatherApi/internal/service/weather"

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

	response, err := h.service.GetWeather(c, city)
	if err != nil {
		c.AbortWithStatusJSON(err.Code, gin.H{"error": err.Message})
		return
	}

	c.JSON(http.StatusOK, response)
}

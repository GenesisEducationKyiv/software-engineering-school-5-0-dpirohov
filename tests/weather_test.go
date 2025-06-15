package tests

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"weatherApi/internal/common/errors"
	"weatherApi/internal/provider"
	"weatherApi/internal/server/routes"
	"weatherApi/internal/service/weather"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestWeatherHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockProvider := &provider.MockProvider{
		Response: &provider.WeatherResponse{
			Temperature: 25.5,
			Humidity:    60,
			Description: "Sunny",
		},
		Err: nil,
	}

	mockFallbackProvider := &provider.MockProvider{
		Response: &provider.WeatherResponse{
			Temperature: 0,
			Humidity:    33,
			Description: "Cloudy",
		},
		Err: nil,
	}

	svc := &weather.WeatherService{
		MainProvider:     mockProvider,
		FallbackProvider: mockFallbackProvider,
	}

	handler := routes.NewWeatherHandler(svc)

	router := gin.Default()
	router.GET("/weather", handler.GetWeather)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "/weather?city=Kyiv", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	var actualResponse provider.WeatherResponse
	err := json.Unmarshal(resp.Body.Bytes(), &actualResponse)
	assert.NoError(t, err)

	assert.Equal(t, 25.5, actualResponse.Temperature)
	assert.Equal(t, 60, actualResponse.Humidity)
	assert.Equal(t, "Sunny", actualResponse.Description)
}

func TestWeatherHandler_MissingCity(t *testing.T) {
	gin.SetMode(gin.TestMode)

	svc := &weather.WeatherService{
		MainProvider:     &provider.MockProvider{},
		FallbackProvider: &provider.MockProvider{},
	}
	handler := routes.NewWeatherHandler(svc)

	router := gin.Default()
	router.GET("/weather", handler.GetWeather)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "/weather", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Body.String(), "city is required")
}

func TestWeatherHandler_FallbackUsed(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mainProvider := &provider.MockProvider{
		Response: nil,
		Err:      errors.New(500, "main failed", nil),
	}

	fallbackProvider := &provider.MockProvider{
		Response: &provider.WeatherResponse{
			Temperature: 15.0,
			Humidity:    80,
			Description: "Cloudy",
		},
		Err: nil,
	}

	svc := &weather.WeatherService{
		MainProvider:     mainProvider,
		FallbackProvider: fallbackProvider,
	}

	handler := routes.NewWeatherHandler(svc)

	router := gin.Default()
	router.GET("/weather", handler.GetWeather)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "/weather?city=London", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	var actualResponse provider.WeatherResponse

	err := json.Unmarshal(resp.Body.Bytes(), &actualResponse)
	assert.NoError(t, err)

	assert.Equal(t, 15.0, actualResponse.Temperature)
	assert.Equal(t, 80, actualResponse.Humidity)
	assert.Equal(t, "Cloudy", actualResponse.Description)
}

func TestCityNotFound(t *testing.T) {
	// Service trusts it's main provider on city search as it's more reliable so we do not fall to second provider if first returns 404
	gin.SetMode(gin.TestMode)

	mainProvider := &provider.MockProvider{
		Response: nil,
		Err:      errors.New(404, "City not found", nil),
	}

	fallbackProvider := &provider.MockProvider{
		Response: &provider.WeatherResponse{
			Temperature: 15.0,
			Humidity:    80,
			Description: "Cloudy",
		},
		Err: nil,
	}

	svc := &weather.WeatherService{
		MainProvider:     mainProvider,
		FallbackProvider: fallbackProvider,
	}

	handler := routes.NewWeatherHandler(svc)

	router := gin.Default()
	router.GET("/weather", handler.GetWeather)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "/weather?city=London", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
	assert.Contains(t, resp.Body.String(), "City not found")
}

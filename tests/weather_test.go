package tests

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"weatherApi/internal/common/utils"
	"weatherApi/internal/dto"

	"weatherApi/internal/provider"
	"weatherApi/internal/server/routes"
	"weatherApi/internal/service/weather"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestWeatherHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockResp := utils.RandomWeatherAPIResponse()

	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(mockResp); err != nil {
			t.Fatalf("failed to encode mock response: %v", err)
		}
	}))
	defer mockAPI.Close()

	mainProvider := provider.NewWeatherApiProvider("test", mockAPI.URL)

	mockFallbackProvider := &provider.MockProvider{
		Response: &dto.WeatherResponse{
			Temperature: 0,
			Humidity:    33,
			Description: "Cloudy",
		},
		Err: nil,
	}

	svc := weather.NewWeatherService(mainProvider, mockFallbackProvider)

	handler := routes.NewWeatherHandler(svc)

	router := gin.Default()
	router.GET("/weather", handler.GetWeather)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "/weather?city=Kyiv", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	var actualResponse dto.WeatherResponse
	err := json.Unmarshal(resp.Body.Bytes(), &actualResponse)
	assert.NoError(t, err)

	assert.Equal(t, mockResp.Current.Temperature, actualResponse.Temperature)
	assert.Equal(t, mockResp.Current.Humidity, actualResponse.Humidity)
	assert.Equal(t, mockResp.Current.Condition.Text, actualResponse.Description)
}

func TestWeatherHandler_MissingCity(t *testing.T) {
	gin.SetMode(gin.TestMode)

	svc := weather.NewWeatherService(&provider.MockProvider{}, &provider.MockProvider{})

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

	mockWeatherAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(map[string]string{
			"error": "Internal server Error",
		}); err != nil {
			t.Fatalf("failed to encode mock response: %v", err)
		}
	}))
	defer mockWeatherAPI.Close()

	mockOpenWeatherMapResp := utils.RandomOpenweatherMapAPIResponse()
	mockOpenWeatherMapAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(mockOpenWeatherMapResp); err != nil {
			t.Fatalf("failed to encode mock response: %v", err)
		}
	}))
	defer mockOpenWeatherMapAPI.Close()

	mainProvider := provider.NewWeatherApiProvider("test", mockWeatherAPI.URL)
	fallbackProvider := provider.NewOpenWeatherApiProvider("test", mockOpenWeatherMapAPI.URL)

	svc := weather.NewWeatherService(mainProvider, fallbackProvider)

	handler := routes.NewWeatherHandler(svc)

	router := gin.Default()
	router.GET("/weather", handler.GetWeather)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "/weather?city=London", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	var actualResponse dto.WeatherResponse

	err := json.Unmarshal(resp.Body.Bytes(), &actualResponse)
	assert.NoError(t, err)

	assert.Equal(t, mockOpenWeatherMapResp.Main.Temperature, actualResponse.Temperature)
	assert.Equal(t, mockOpenWeatherMapResp.Main.Humidity, actualResponse.Humidity)
	assert.Equal(t, mockOpenWeatherMapResp.Weather[0].Description, actualResponse.Description)
}

func TestCityNotFound(t *testing.T) {
	// Service trusts it's main provider on city search as it's more reliable so we do not fall to second provider if first returns 404
	gin.SetMode(gin.TestMode)
	mockWeatherAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(map[string]string{
			"error": "City not found",
		}); err != nil {
			t.Fatalf("failed to encode mock response: %v", err)
		}
	}))
	defer mockWeatherAPI.Close()

	mockOpenWeatherMapResp := utils.RandomOpenweatherMapAPIResponse()
	mockOpenWeatherMapAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(mockOpenWeatherMapResp); err != nil {
			t.Fatalf("failed to encode mock response: %v", err)
		}
	}))
	defer mockOpenWeatherMapAPI.Close()

	mainProvider := provider.NewWeatherApiProvider("test", mockWeatherAPI.URL)
	fallbackProvider := provider.NewOpenWeatherApiProvider("test", mockOpenWeatherMapAPI.URL)

	svc := weather.NewWeatherService(mainProvider, fallbackProvider)

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

func TestWeatherHandler_RealProviderWithMockedAPI(t *testing.T) {
	mockResp := utils.RandomOpenweatherMapAPIResponse()

	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(mockResp); err != nil {
			t.Fatalf("failed to encode mock response: %v", err)
		}
	}))
	defer mockAPI.Close()

	prov := provider.NewOpenWeatherApiProvider("test", mockAPI.URL)

	svc := weather.NewWeatherService(prov)
	handler := routes.NewWeatherHandler(svc)

	router := gin.Default()
	router.GET("/weather", handler.GetWeather)

	req := httptest.NewRequest(http.MethodGet, "/weather?city=Kyiv", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var data dto.WeatherResponse
	err := json.Unmarshal(resp.Body.Bytes(), &data)
	assert.NoError(t, err)
	assert.Equal(t, mockResp.Main.Temperature, data.Temperature)
	assert.Equal(t, mockResp.Main.Humidity, data.Humidity)
	assert.Equal(t, mockResp.Weather[0].Description, data.Description)
}

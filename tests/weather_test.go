package tests

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"weatherApi/internal/logger"

	"github.com/stretchr/testify/require"

	"weatherApi/internal/common/utils"
	"weatherApi/internal/dto"

	"weatherApi/internal/provider"
	cacheRepo "weatherApi/internal/repository/weather"
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
	log := logger.NewNoOpLogger()

	mainProvider := provider.NewWeatherApiProvider(log, "test", mockAPI.URL)

	mockFallbackProvider := &provider.MockProvider{
		Response: &dto.WeatherResponse{
			Temperature: 0,
			Humidity:    33,
			Description: "Cloudy",
		},
		Err: nil,
	}
	mockCacheRepo := cacheRepo.NewMockCacheRepo()

	svc := weather.NewWeatherService(log, mockCacheRepo, mainProvider, mockFallbackProvider)

	handler := routes.NewWeatherHandler(log, svc)

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
	mockCacheRepo := cacheRepo.NewMockCacheRepo()
	log := logger.NewNoOpLogger()

	svc := weather.NewWeatherService(log, mockCacheRepo, &provider.MockProvider{}, &provider.MockProvider{})

	handler := routes.NewWeatherHandler(log, svc)

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
	log := logger.NewNoOpLogger()

	mainProvider := provider.NewWeatherApiProvider(log, "test", mockWeatherAPI.URL)
	fallbackProvider := provider.NewOpenWeatherApiProvider(log, "test", mockOpenWeatherMapAPI.URL)
	mockCacheRepo := cacheRepo.NewMockCacheRepo()

	svc := weather.NewWeatherService(log, mockCacheRepo, mainProvider, fallbackProvider)

	handler := routes.NewWeatherHandler(log, svc)

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
	log := logger.NewNoOpLogger()

	mainProvider := provider.NewWeatherApiProvider(log, "test", mockWeatherAPI.URL)
	fallbackProvider := provider.NewOpenWeatherApiProvider(log, "test", mockOpenWeatherMapAPI.URL)
	mockCacheRepo := cacheRepo.NewMockCacheRepo()

	svc := weather.NewWeatherService(log, mockCacheRepo, mainProvider, fallbackProvider)

	handler := routes.NewWeatherHandler(log, svc)

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
	log := logger.NewNoOpLogger()

	prov := provider.NewOpenWeatherApiProvider(log, "test", mockAPI.URL)
	mockCacheRepo := cacheRepo.NewMockCacheRepo()

	svc := weather.NewWeatherService(log, mockCacheRepo, prov)
	handler := routes.NewWeatherHandler(log, svc)

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

func TestWeatherService_Cache_HTTP(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mockAPIResp := utils.RandomWeatherAPIResponse()
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(mockAPIResp); err != nil {
			t.Fatalf("failed to encode mock response: %v", err)
		}
	}))
	defer mockAPI.Close()
	log := logger.NewNoOpLogger()

	mainProvider := provider.NewWeatherApiProvider(log, "test", mockAPI.URL)
	mockRepo := cacheRepo.NewMockCacheRepo()
	svc := weather.NewWeatherService(log, mockRepo, mainProvider)

	handler := routes.NewWeatherHandler(log, svc)
	router := gin.Default()
	router.GET("/weather", handler.GetWeather)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	city := "TestCity"

	req1, _ := http.NewRequestWithContext(ctx, http.MethodGet, "/weather?city="+city, nil)
	resp1 := httptest.NewRecorder()
	router.ServeHTTP(resp1, req1)

	assert.Equal(t, http.StatusOK, resp1.Code)

	var data1 dto.WeatherResponse
	err := json.Unmarshal(resp1.Body.Bytes(), &data1)
	assert.NoError(t, err)

	// Check response is received from cache
	req2, _ := http.NewRequestWithContext(ctx, http.MethodGet, "/weather?city="+city, nil)
	resp2 := httptest.NewRecorder()
	router.ServeHTTP(resp2, req2)

	assert.Equal(t, http.StatusOK, resp2.Code)

	var data2 dto.WeatherResponse
	err = json.Unmarshal(resp2.Body.Bytes(), &data2)
	assert.NoError(t, err)
	assert.Equal(t, data1.Temperature, data2.Temperature)
	assert.Equal(t, data1.Humidity, data2.Humidity)
	assert.Equal(t, data1.Description, data2.Description)
}

func TestWeatherService_WaitForLockedCache(t *testing.T) {
	gin.SetMode(gin.TestMode)
	city := "TestCity"
	expectedResponse := utils.RandomWeatherAPIResponse()

	mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(utils.RandomWeatherAPIResponse()); err != nil {
			t.Fatalf("failed to encode mock data: %v", err)
		}
	}))
	defer mockAPI.Close()
	log := logger.NewNoOpLogger()

	mainProvider := provider.NewWeatherApiProvider(log, "test", mockAPI.URL)
	mockRepo := cacheRepo.NewMockCacheRepo()
	svc := weather.NewWeatherService(log, mockRepo, mainProvider)

	handler := routes.NewWeatherHandler(log, svc)
	router := gin.Default()
	router.GET("/weather", handler.GetWeather)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// emulate lock is acquired by another process
	locked, errLock := mockRepo.AcquireLock(ctx, city)
	assert.True(t, locked)
	assert.NoError(t, errLock)

	go func() {
		time.Sleep(300 * time.Millisecond)

		err := mockRepo.Set(ctx, city, &dto.WeatherResponse{
			Temperature: expectedResponse.Current.Temperature,
			Humidity:    expectedResponse.Current.Humidity,
			Description: expectedResponse.Current.Condition.Text,
		})
		if err != nil {
			t.Errorf("failed to set cache in goroutine: %v", err)
		}
		err = mockRepo.ReleaseLock(ctx, city)
		if err != nil {
			t.Errorf("failed to release lock in goroutine: %v", err)
		}
	}()

	start := time.Now()
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "/weather?city="+city, nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	duration := time.Since(start)

	assert.Equal(t, http.StatusOK, resp.Code)

	// Check response is received after lock is released
	var data dto.WeatherResponse
	err := json.Unmarshal(resp.Body.Bytes(), &data)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse.Current.Temperature, data.Temperature)
	assert.Equal(t, expectedResponse.Current.Humidity, data.Humidity)
	assert.Equal(t, expectedResponse.Current.Condition.Text, data.Description)
	assert.GreaterOrEqual(t, duration.Milliseconds(), int64(290), "Request should wait for lock release")
}

func TestWeatherService_ProviderCallCount(t *testing.T) {
	gin.SetMode(gin.TestMode)
	city := "TestCity"
	mockResp := utils.RandomWeatherAPIResponse()

	mockProv := &provider.MockProvider{
		Response: &dto.WeatherResponse{
			Temperature: mockResp.Current.Temperature,
			Humidity:    mockResp.Current.Humidity,
			Description: mockResp.Current.Condition.Text,
		},
		Err: nil,
	}
	mockRepo := cacheRepo.NewMockCacheRepo()
	log := logger.NewNoOpLogger()

	svc := weather.NewWeatherService(log, mockRepo, mockProv)

	handler := routes.NewWeatherHandler(log, svc)
	router := gin.Default()
	router.GET("/weather", handler.GetWeather)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req1, _ := http.NewRequestWithContext(ctx, http.MethodGet, "/weather?city="+city, nil)
	resp1 := httptest.NewRecorder()
	router.ServeHTTP(resp1, req1)

	assert.Equal(t, http.StatusOK, resp1.Code)
	require.Equal(t, 1, mockProv.GetWeatherCallCount)

	// expect get weather was not called as data is cached
	req2, _ := http.NewRequestWithContext(ctx, http.MethodGet, "/weather?city="+city, nil)
	resp2 := httptest.NewRecorder()
	router.ServeHTTP(resp2, req2)

	assert.Equal(t, http.StatusOK, resp1.Code)
	require.NotNil(t, resp2)
	require.Equal(t, 1, mockProv.GetWeatherCallCount)

	// expect get weather was not called as data is cached
	req3, _ := http.NewRequestWithContext(ctx, http.MethodGet, "/weather?city="+city, nil)
	resp3 := httptest.NewRecorder()
	router.ServeHTTP(resp3, req3)

	assert.Equal(t, http.StatusOK, resp1.Code)
	require.NotNil(t, resp3)
	require.Equal(t, 1, mockProv.GetWeatherCallCount)
}

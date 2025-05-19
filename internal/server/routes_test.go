package server

import (
	"weatherApi/internal/server/routes"
	"weatherApi/internal/service/weather"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGetWeatherHandler(t *testing.T) {
	r := gin.New()
	weatherHandler := routes.NewWeatherHandler(weather.NewWeatherService())
	r.GET("/weather", weatherHandler.GetWeather)
	
	// Create a test HTTP request
	req, err := http.NewRequest("GET", "/weather", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Serve the HTTP request
	r.ServeHTTP(rr, req)
	
	// Check the status code
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

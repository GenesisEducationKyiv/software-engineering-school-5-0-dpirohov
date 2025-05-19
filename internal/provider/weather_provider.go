package provider

import (
	"weatherApi/internal/common/errors"
	"net/http"
)

type WeatherResponse struct {
	Temperature float64 `json:"temperature"`
	Humidity    int     `json:"humidity"`
	Description string  `json:"description"`
}

type WeatherProviderInterface interface {
	GetWeather(city string) (*WeatherResponse, *errors.AppError)
	checkApiResponse(response *http.Response) *errors.AppError
	handleInternalError(err error) *errors.AppError
}

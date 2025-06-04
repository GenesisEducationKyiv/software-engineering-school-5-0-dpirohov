package provider

import (
	"net/http"

	"weatherApi/internal/common/errors"
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

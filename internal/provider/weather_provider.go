package provider

import (
	"weatherApi/internal/common/errors"
	"weatherApi/internal/dto"
)

type WeatherProviderInterface interface {
	SetNext(next WeatherProviderInterface)
	Next(city string) (*dto.WeatherResponse, *errors.AppError)
	GetWeather(city string) (*dto.WeatherResponse, *errors.AppError)
}

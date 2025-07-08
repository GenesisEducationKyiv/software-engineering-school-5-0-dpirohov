package provider

import (
	"log"
	"weatherApi/internal/common/errors"
	"weatherApi/internal/dto"
	serviceErrors "weatherApi/internal/service/weather/errors"
)

type WeatherProviderInterface interface {
	SetNext(next WeatherProviderInterface)
	GetWeather(city string) (*dto.WeatherResponse, *errors.AppError)
	Name() string
}

func TryNext(current WeatherProviderInterface, next WeatherProviderInterface, city string, err error) (*dto.WeatherResponse, *errors.AppError) {
	log.Printf("%s: error: %v — trying next provider...", current.Name(), err)

	if next != nil {
		return next.GetWeather(city)
	}

	log.Printf("%s: no next provider available", current.Name())
	return nil, serviceErrors.ErrInternalServerError
}

package provider

import (
	"log"
	"reflect"
	"weatherApi/internal/common/errors"
	"weatherApi/internal/dto"
	serviceErrors "weatherApi/internal/service/weather/errors"
)

type WeatherProviderInterface interface {
	SetNext(next WeatherProviderInterface)
	GetWeather(city string) (*dto.WeatherResponse, *errors.AppError)
}

func TryNext(current WeatherProviderInterface, next WeatherProviderInterface, city string, err error) (*dto.WeatherResponse, *errors.AppError) {
	providerName := reflect.TypeOf(current).Elem().Name()
	log.Printf("%s: error: %v — trying next provider...", providerName, err)

	if next != nil {
		return next.GetWeather(city)
	}

	log.Printf("%s: no next provider available", providerName)
	return nil, serviceErrors.ErrInternalServerError
}

package weather

import (
	"log"

	"weatherApi/internal/common/errors"
	"weatherApi/internal/provider"
)

type WeatherService struct {
	MainProvider     provider.WeatherProviderInterface
	FallbackProvider provider.WeatherProviderInterface
}

func NewWeatherService(mainProvider, fallbackProvider provider.WeatherProviderInterface) *WeatherService {
	return &WeatherService{
		MainProvider:     mainProvider,
		FallbackProvider: fallbackProvider,
	}
}

func (service *WeatherService) GetWeather(
	city string,
) (*provider.WeatherResponse, *errors.AppError) {
	response, err := service.MainProvider.GetWeather(city)
	if err == nil {
		return response, nil
	} else if err.Code != 500 {
		return nil, err
	}

	log.Printf("Main provider error: %s; trying fallback", err.Message)

	fallbackResponse, fallbackErr := service.FallbackProvider.GetWeather(city)

	if fallbackErr != nil {
		log.Printf("Fallback provider error: %s", fallbackErr.Message)
		return nil, fallbackErr
	}

	return fallbackResponse, nil
}

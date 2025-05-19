package weather

import (
	"weatherApi/internal/common/errors"
	"weatherApi/internal/provider"
	"log"
	"os"
)

type WeatherService struct {
	MainProvider     provider.WeatherProviderInterface
	FallbackProvider provider.WeatherProviderInterface
}

func NewWeatherService() *WeatherService {
	return &WeatherService{
		MainProvider:     provider.NewOpenWeatherApiProvider(os.Getenv("OPENWEATHER__API_KEY")),
		FallbackProvider: provider.NewWeatherApiProvider(os.Getenv("WEATHER_API_API_KEY")),
	}
}

func (service *WeatherService) GetWeather(city string) (*provider.WeatherResponse, *errors.AppError) {


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

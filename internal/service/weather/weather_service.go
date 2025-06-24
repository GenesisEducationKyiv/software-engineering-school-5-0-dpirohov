package weather

import (
	"weatherApi/internal/dto"

	"weatherApi/internal/common/errors"
	"weatherApi/internal/provider"
)

type Service struct {
	provider provider.WeatherProviderInterface
}

func NewWeatherService(providers ...provider.WeatherProviderInterface) *Service {
	for i := 0; i < len(providers)-1; i++ {
		providers[i].SetNext(providers[i+1])
	}
	return &Service{provider: providers[0]}
}

func (service *Service) GetWeather(
	city string,
) (*dto.WeatherResponse, *errors.AppError) {
	return service.provider.GetWeather(city)
}

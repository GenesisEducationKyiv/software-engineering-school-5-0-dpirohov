package provider

import (
	"context"
	"weatherApi/internal/common/errors"
	"weatherApi/internal/dto"

	"github.com/rs/zerolog"

	serviceErrors "weatherApi/internal/service/weather/errors"
)

type WeatherProviderInterface interface {
	SetNext(next WeatherProviderInterface)
	GetWeather(ctx context.Context, city string) (*dto.WeatherResponse, *errors.AppError)
	Name() string
}

func TryNext(log *zerolog.Logger, ctx context.Context, current WeatherProviderInterface, next WeatherProviderInterface, city string, err error) (*dto.WeatherResponse, *errors.AppError) {
	log.Error().Err(err).Msgf("%s: Provider failed", current.Name())

	if next != nil {
		return next.GetWeather(ctx, city)
	}

	log.Error().Msgf("%s: no next provider available", current.Name())
	return nil, serviceErrors.ErrInternalServerError
}

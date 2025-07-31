package weather

import (
	"context"
	"errors"
	"weatherApi/internal/dto"
	"weatherApi/internal/logger"
	"weatherApi/internal/repository/weather"

	appErrors "weatherApi/internal/common/errors"
	serviceErrors "weatherApi/internal/service/weather/errors"

	"weatherApi/internal/provider"
)

type Service struct {
	provider  provider.WeatherProviderInterface
	cacheRepo weather.CacheRepoInterface
}

func NewWeatherService(cacheRepo weather.CacheRepoInterface, providers ...provider.WeatherProviderInterface) *Service {
	if len(providers) == 0 {
		panic("At least one provider required!")
	}
	for i := 0; i < len(providers)-1; i++ {
		providers[i].SetNext(providers[i+1])
	}
	return &Service{provider: providers[0], cacheRepo: cacheRepo}
}

func (service *Service) GetWeather(
	ctx context.Context,
	city string,
) (*dto.WeatherResponse, *appErrors.AppError) {
	log := logger.FromContext(ctx)

	resp, err := service.cacheRepo.Get(ctx, city)
	if err != nil && !errors.Is(err, weather.ErrCacheIsEmpty) {
		log.Error().Err(err).Msg("Redis error, caching is skipped!")
		return service.provider.GetWeather(ctx, city)
	}
	if resp != nil {
		return resp, nil
	}

	locked, err := service.cacheRepo.AcquireLock(ctx, city)
	if err != nil {
		log.Error().Err(err).Msg("Redis failed to acquire lock")
	}
	if !locked {
		response, err := service.cacheRepo.WaitForUnlock(ctx, city)
		if err != nil {
			return nil, serviceErrors.ErrInternalServerError
		}
		if response != nil {
			return response, nil
		}
	} else {
		defer func(cacheRepo weather.CacheRepoInterface, ctx context.Context, city string) {
			err := cacheRepo.ReleaseLock(ctx, city)
			if err != nil {
				log.Error().Err(err).Msg("Failed to release lock")
			}
		}(service.cacheRepo, ctx, city)
	}

	result, appErr := service.provider.GetWeather(ctx, city)
	if appErr != nil {
		return nil, appErr
	}

	_ = service.cacheRepo.Set(ctx, city, result)
	return result, nil
}

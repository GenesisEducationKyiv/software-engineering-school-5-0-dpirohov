package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"weatherApi/internal/dto"
	"weatherApi/internal/logger"

	"weatherApi/internal/common/errors"
	serviceErrors "weatherApi/internal/service/weather/errors"
)

var _ WeatherProviderInterface = (*WeatherApiProvider)(nil)

type WeatherApiProvider struct {
	log    *logger.Logger
	next   WeatherProviderInterface
	apiKey string
	url    string
}

func NewWeatherApiProvider(log *logger.Logger, apikey, url string) *WeatherApiProvider {
	return &WeatherApiProvider{
		log:    log,
		apiKey: apikey,
		url:    url,
	}
}

func (w *WeatherApiProvider) Name() string {
	return "WeatherApi"
}

func (w *WeatherApiProvider) SetNext(next WeatherProviderInterface) {
	w.next = next
}

func (w *WeatherApiProvider) Next(ctx context.Context, city string) (*dto.WeatherResponse, *errors.AppError) {
	log := w.log.FromContext(ctx)
	if w.next != nil {
		return w.next.GetWeather(ctx, city)
	}
	log.Error().Msg("WeatherApiProvider: no providers left in chain!")
	return nil, serviceErrors.ErrInternalServerError
}

func (w *WeatherApiProvider) GetWeather(ctx context.Context, city string) (*dto.WeatherResponse, *errors.AppError) {
	log := w.log.FromContext(ctx)

	var weatherResponse dto.WeatherAPIResponse
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		fmt.Sprintf("%s?key=%s&q=%s&aqi=no", w.url, w.apiKey, city),
		nil,
	)
	if err != nil {
		return TryNext(log, ctx, w, w.next, city, err)
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return TryNext(log, ctx, w, w.next, city, err)
	}

	defer func() {
		if err := response.Body.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close response body")
		}
	}()

	if badResponse := w.checkApiResponse(response); badResponse != nil {
		if badResponse.Code == 500 {
			return TryNext(log, ctx, w, w.next, city, fmt.Errorf("bad API response: %v", badResponse.Message))
		}
		return nil, badResponse
	}

	if err := json.NewDecoder(response.Body).Decode(&weatherResponse); err != nil {
		return TryNext(log, ctx, w, w.next, city, fmt.Errorf("failed to decode response: %w", err))
	}

	return &dto.WeatherResponse{
		Temperature: weatherResponse.Current.Temperature,
		Humidity:    weatherResponse.Current.Humidity,
		Description: weatherResponse.Current.Condition.Text,
	}, nil
}

func (w *WeatherApiProvider) checkApiResponse(response *http.Response) *errors.AppError {
	switch response.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusNotFound:
		return serviceErrors.ErrCityNotFound
	default:
		return serviceErrors.ErrInternalServerError
	}
}

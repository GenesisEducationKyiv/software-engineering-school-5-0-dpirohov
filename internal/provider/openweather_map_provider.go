package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"weatherApi/internal/dto"
	"weatherApi/internal/logger"

	"github.com/rs/zerolog/log"

	"weatherApi/internal/common/errors"
	serviceErrors "weatherApi/internal/service/weather/errors"
)

var _ WeatherProviderInterface = (*WeatherApiProvider)(nil)

type OpenWeatherMapApiProvider struct {
	next   WeatherProviderInterface
	apiKey string
	url    string
}

func NewOpenWeatherApiProvider(apikey, url string) *OpenWeatherMapApiProvider {
	return &OpenWeatherMapApiProvider{
		apiKey: apikey,
		url:    url,
	}
}

func (w *OpenWeatherMapApiProvider) Name() string {
	return "OpenWeatherMap"
}

func (w *OpenWeatherMapApiProvider) SetNext(next WeatherProviderInterface) {
	w.next = next
}

func (w *OpenWeatherMapApiProvider) Next(ctx context.Context, city string) (*dto.WeatherResponse, *errors.AppError) {
	log := logger.FromContext(ctx)
	if w.next != nil {
		return w.next.GetWeather(ctx, city)
	}
	log.Error().Msg("OpenWeatherMapApiProvider: no providers left in chain!")
	return nil, serviceErrors.ErrInternalServerError
}

func (w *OpenWeatherMapApiProvider) GetWeather(ctx context.Context, city string) (*dto.WeatherResponse, *errors.AppError) {
	var openWeatherMapResponse dto.OpenweatherMapAPIResponse
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		fmt.Sprintf("%s?q=%s&APPID=%s&units=metric", w.url, city, w.apiKey),
		nil,
	)
	if err != nil {
		return TryNext(ctx, w, w.next, city, fmt.Errorf("request creation failed: %w", err))
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return TryNext(ctx, w, w.next, city, fmt.Errorf("HTTP request failed: %w", err))
	}

	defer func() {
		if err := response.Body.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close response body")
		}
	}()

	if badResponse := w.checkApiResponse(response); badResponse != nil {
		if badResponse.Code == 500 {
			return TryNext(ctx, w, w.next, city, fmt.Errorf("bad API response: %w", err))
		}
		return nil, badResponse
	}

	if err := json.NewDecoder(response.Body).Decode(&openWeatherMapResponse); err != nil {
		return TryNext(ctx, w, w.next, city, fmt.Errorf("failed to decode response: %w", err))
	}

	var description string

	if len(openWeatherMapResponse.Weather) > 0 {
		description = openWeatherMapResponse.Weather[0].Description
	}

	return &dto.WeatherResponse{
		Temperature: openWeatherMapResponse.Main.Temperature,
		Humidity:    openWeatherMapResponse.Main.Humidity,
		Description: description,
	}, nil
}

func (w *OpenWeatherMapApiProvider) checkApiResponse(response *http.Response) *errors.AppError {
	switch response.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusNotFound:
		return serviceErrors.ErrCityNotFound
	default:
		return serviceErrors.ErrInternalServerError
	}
}

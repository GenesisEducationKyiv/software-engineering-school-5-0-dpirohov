package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"weatherApi/internal/dto"

	"weatherApi/internal/common/errors"
	serviceErrors "weatherApi/internal/service/weather/errors"
)

var _ WeatherProviderInterface = (*WeatherApiProvider)(nil)

type WeatherApiProvider struct {
	next   WeatherProviderInterface
	apiKey string
	url    string
}

func NewWeatherApiProvider(apikey, url string) *WeatherApiProvider {
	return &WeatherApiProvider{
		apiKey: apikey,
		url:    url,
	}
}
func (w *WeatherApiProvider) SetNext(next WeatherProviderInterface) {
	w.next = next
}

func (w *WeatherApiProvider) Next(city string) (*dto.WeatherResponse, *errors.AppError) {
	if w.next != nil {
		return w.next.GetWeather(city)
	}
	log.Printf("WeatherApiProvider: no providers left in chain!")
	return nil, serviceErrors.ErrInternalServerError
}

func (w *WeatherApiProvider) GetWeather(city string) (*dto.WeatherResponse, *errors.AppError) {
	var weatherResponse dto.WeatherAPIResponse
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		fmt.Sprintf("%s?key=%s&q=%s&aqi=no", w.url, w.apiKey, city),
		nil,
	)
	if err != nil {
		return TryNext(w, w.next, city, err)
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return TryNext(w, w.next, city, err)
	}

	defer func() {
		if err := response.Body.Close(); err != nil {
			log.Printf("failed to close response body: %v", err)
		}
	}()

	if badResponse := w.checkApiResponse(response); badResponse != nil {
		if badResponse.Code == 500 {
			return TryNext(w, w.next, city, fmt.Errorf("bad API response: %v", badResponse.Message))
		}
		return nil, badResponse
	}

	if err := json.NewDecoder(response.Body).Decode(&weatherResponse); err != nil {
		return TryNext(w, w.next, city, fmt.Errorf("failed to decode response: %w", err))
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

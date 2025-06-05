package provider

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"weatherApi/internal/common/errors"
	serviceErrors "weatherApi/internal/service/weather/errors"
)

type weatherAPIResponse struct {
	Location struct {
		Name    string `json:"name"`
		Region  string `json:"region"`
		Country string `json:"country"`
	} `json:"location"`

	Current struct {
		Temperature float64 `json:"temp_c"`
		Humidity    int     `json:"humidity"`
		Condition   struct {
			Text string `json:"text"`
		} `json:"condition"`
	} `json:"current"`
}

type WeatherApiProvider struct {
	apiKey string
	url    string
}

func NewWeatherApiProvider(apikey string) WeatherProviderInterface {
	return &WeatherApiProvider{apiKey: apikey, url: "http://api.weatherapi.com/v1/current.json"}
}

func (w *WeatherApiProvider) GetWeather(city string) (*WeatherResponse, *errors.AppError) {
	var weatherResponse weatherAPIResponse
	response, err := http.Get(fmt.Sprintf("%s?key=%s&q=%s&aqi=no", w.url, w.apiKey, city))
	if err != nil {
		return nil, w.handleInternalError(err)
	}

	if badResponse := w.checkApiResponse(response); badResponse != nil {
		return nil, badResponse
	}

	defer func() {
		if err := response.Body.Close(); err != nil {
			log.Printf("failed to close response body: %v", err)
		}
	}()

	if err := json.NewDecoder(response.Body).Decode(&weatherResponse); err != nil {
		return nil, w.handleInternalError(err)
	}

	return &WeatherResponse{
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

func (w *WeatherApiProvider) handleInternalError(err error) *errors.AppError {
	log.Printf("WeatherApiProvider HTTP request failed: %v", err)
	return errors.New(500, fmt.Errorf("internal server error: %w", err).Error(), err)
}

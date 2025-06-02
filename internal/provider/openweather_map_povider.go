package provider

import (
	"encoding/json"
	"fmt"
	"log"
	"weatherApi/internal/common/errors"
	serviceErrors "weatherApi/internal/service/weather/errors"

	"net/http"
)

type openweatherMapAPIResponse struct {
	Main struct {
		Temperature float64 `json:"temp"`
		FeelsLike   float64 `json:"feels_like"`
		Humidity    int     `json:"humidity"`
	} `json:"main"`

	Weather []struct {
		ID          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
}

type OpenWeatherMapApiProvider struct {
	apiKey string
	url    string
}

func NewOpenWeatherApiProvider(apikey string) WeatherProviderInterface {
	return &OpenWeatherMapApiProvider{apiKey: apikey, url: "http://api.openweathermap.org/data/2.5/weather"}
}

func (w *OpenWeatherMapApiProvider) GetWeather(city string) (*WeatherResponse, *errors.AppError) {
	var openWeatherMapResponse openweatherMapAPIResponse
	response, error := http.Get(fmt.Sprintf("%s?q=%s&APPID=%s&units=metric", w.url, city, w.apiKey))

	if error != nil {
		return nil, w.handleInternalError(error)
	}

	if badResponse := w.checkApiResponse(response); badResponse != nil {
		return nil, badResponse
	}

	defer response.Body.Close()

	if err := json.NewDecoder(response.Body).Decode(&openWeatherMapResponse); err != nil {
		return nil, w.handleInternalError(error)
	}

	var description string

	if len(openWeatherMapResponse.Weather) > 0 {
		description = openWeatherMapResponse.Weather[0].Description
	}

	return &WeatherResponse{
		Temperature: openWeatherMapResponse.Main.Temperature,
		Humidity:    openWeatherMapResponse.Main.Humidity,
		Description: description,
	}, nil
}

func (w *OpenWeatherMapApiProvider) checkApiResponse(response *http.Response) *errors.AppError {
	switch response.StatusCode {
	case 200:
		return nil
	case 404:
		return serviceErrors.CityNotFoundError
	default:
		return serviceErrors.InternalServerError
	}
}

func (w *OpenWeatherMapApiProvider) handleInternalError(err error) *errors.AppError {
	log.Printf("OpenWeatherMapApiProvider HTTP request failed: %v", err)
	return errors.New(500, fmt.Errorf("internal server error: %w", err).Error(), err)
}

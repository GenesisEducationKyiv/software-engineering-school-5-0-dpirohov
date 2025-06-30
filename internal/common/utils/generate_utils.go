package utils

import (
	"math/rand"
	"time"
	"weatherApi/internal/dto"
)

func RandomOpenweatherMapAPIResponse() dto.OpenweatherMapAPIResponse {
	// #nosec G404 -- not used for cryptographic purposes
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	weatherConditions := []struct {
		ID          int
		Main        string
		Description string
		Icon        string
	}{
		{800, "Clear", "Sunny", "01d"},
		{801, "Clouds", "Few clouds", "02d"},
		{802, "Clouds", "Scattered clouds", "03d"},
		{803, "Clouds", "Broken clouds", "04d"},
		{500, "Rain", "Light rain", "10d"},
		{600, "Snow", "Light snow", "13d"},
		{701, "Mist", "Mist", "50d"},
	}

	condition := weatherConditions[r.Intn(len(weatherConditions))]

	temp := r.Float64()*40 - 10
	feelsLike := temp + (r.Float64()*4 - 2)
	humidity := r.Intn(100)

	resp := dto.OpenweatherMapAPIResponse{}
	resp.Main.Temperature = temp
	resp.Main.FeelsLike = feelsLike
	resp.Main.Humidity = humidity

	resp.Weather = []struct {
		ID          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	}{
		{
			ID:          condition.ID,
			Main:        condition.Main,
			Description: condition.Description,
			Icon:        condition.Icon,
		},
	}

	return resp
}

func RandomWeatherAPIResponse() dto.WeatherAPIResponse {
	// #nosec G404 -- not used for cryptographic purposes
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	cities := []string{"Kyiv", "London", "Berlin", "Tokyo", "New York"}
	countries := []string{"UA", "UK", "DE", "JP", "US"}
	regions := []string{"Region1", "Region2", "Region3", "Region4"}

	descriptions := []string{"Sunny", "Cloudy", "Rainy", "Snowy", "Windy", "Foggy"}

	resp := dto.WeatherAPIResponse{}

	resp.Location.Name = cities[r.Intn(len(cities))]
	resp.Location.Region = regions[r.Intn(len(regions))]
	resp.Location.Country = countries[r.Intn(len(countries))]

	resp.Current.Temperature = r.Float64()*40 - 10
	resp.Current.Humidity = r.Intn(100)
	resp.Current.Condition.Text = descriptions[r.Intn(len(descriptions))]

	return resp
}

package provider

import (
	"net/http"

	"weatherApi/internal/common/errors"
)

type MockProvider struct {
	Response *WeatherResponse
	Err      *errors.AppError
}

func (m *MockProvider) GetWeather(city string) (*WeatherResponse, *errors.AppError) {
	return m.Response, m.Err
}

func (m *MockProvider) checkApiResponse(_ *http.Response) *errors.AppError {
	return nil
}

func (m *MockProvider) handleInternalError(_ error) *errors.AppError {
	return errors.New(500, "internal", nil)
}

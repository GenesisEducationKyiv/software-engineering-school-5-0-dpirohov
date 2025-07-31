package provider

import (
	"context"
	"weatherApi/internal/dto"
	serviceErrors "weatherApi/internal/service/weather/errors"

	"weatherApi/internal/common/errors"
)

type MockProvider struct {
	next                WeatherProviderInterface
	Response            *dto.WeatherResponse
	Err                 *errors.AppError
	GetWeatherCallCount int
}

func (m *MockProvider) GetWeather(ctx context.Context, city string) (*dto.WeatherResponse, *errors.AppError) {
	m.GetWeatherCallCount++
	if m.Err != nil {
		if m.Err.Code == 500 && m.next != nil {
			return m.Next(ctx, city)
		}
		return nil, m.Err
	}
	return m.Response, nil
}

func (m *MockProvider) Name() string {
	return "MockProvider"
}

func (m *MockProvider) SetNext(next WeatherProviderInterface) {
	m.next = next
}

func (m *MockProvider) Next(ctx context.Context, city string) (*dto.WeatherResponse, *errors.AppError) {
	if m.next != nil {
		return m.next.GetWeather(ctx, city)
	}
	return nil, serviceErrors.ErrInternalServerError
}

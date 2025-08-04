package provider

import (
	"weatherApi/internal/dto"
)

type MockSMTPClient struct {
	SentConfirmations []dto.ConfirmationEmailTask
	SentWeatherData   []dto.WeatherResponse
	SentUserData      []dto.UserData
}

func (m *MockSMTPClient) SendConfirmationToken(email, token, city string) error {
	m.SentConfirmations = append(m.SentConfirmations, dto.ConfirmationEmailTask{
		Email: email,
		Token: token,
		City:  city,
	})
	return nil
}

func (m *MockSMTPClient) SendSubscriptionWeatherData(data *dto.WeatherResponse, user *dto.UserData) error {
	m.SentWeatherData = append(m.SentWeatherData, *data)
	m.SentUserData = append(m.SentUserData, *user)
	return nil
}

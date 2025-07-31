package tests

import (
	"context"
	"encoding/json"
	"testing"
	"weatherApi/internal/broker"
	"weatherApi/internal/common/utils"
	"weatherApi/internal/dto"
	"weatherApi/internal/provider"
	"weatherApi/internal/worker"

	"github.com/stretchr/testify/assert"
)

func TestStartConfirmationWorker(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mockSubscriber := broker.NewMockEventSubscriber()
	mockSMTP := &provider.MockSMTPClient{}

	err := worker.StartConfirmationWorker(ctx, mockSubscriber, mockSMTP)
	assert.NoError(t, err)

	task := dto.ConfirmationEmailTask{
		Email: "test@example.com",
		Token: "abc123",
		City:  "Kyiv",
	}
	data, _ := json.Marshal(task)

	err = mockSubscriber.SimulateMessage(ctx, broker.SubscriptionConfirmationTasks, data)
	assert.NoError(t, err)
	assert.Len(t, mockSMTP.SentConfirmations, 1)
	assert.Equal(t, task.Email, mockSMTP.SentConfirmations[0].Email)
	assert.Equal(t, task.City, mockSMTP.SentConfirmations[0].City)
	assert.Equal(t, task.Token, mockSMTP.SentConfirmations[0].Token)
}

func TestStartSubscriptionWorker(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mockSubscriber := broker.NewMockEventSubscriber()
	mockSMTP := &provider.MockSMTPClient{}

	err := worker.StartSubscriptionWorker(ctx, mockSubscriber, mockSMTP)
	assert.NoError(t, err)

	randomResponse := utils.RandomWeatherAPIResponse()

	task := dto.WeatherSubData{
		Weather: dto.WeatherResponse{
			Temperature: randomResponse.Current.Temperature,
			Humidity:    randomResponse.Current.Humidity,
			Description: randomResponse.Current.Condition.Text,
		},
		Users: []dto.UserData{
			{Email: "user1@example.com", Token: "123"},
			{Email: "user2@example.com", Token: "321"},
		},
	}
	data, _ := json.Marshal(task)

	err = mockSubscriber.SimulateMessage(ctx, broker.SendSubscriptionWeatherData, data)
	assert.NoError(t, err)
	assert.Len(t, mockSMTP.SentWeatherData, 2)
	assert.Equal(t, task.Users[0], mockSMTP.SentUserData[1])
	assert.Equal(t, task.Weather, mockSMTP.SentWeatherData[0])

	assert.Equal(t, task.Users[1], mockSMTP.SentUserData[0])
	assert.Equal(t, task.Weather, mockSMTP.SentWeatherData[1])

}

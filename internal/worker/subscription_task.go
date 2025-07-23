package worker

import (
	"context"
	"encoding/json"
	"log"
	"weatherApi/internal/broker"
	"weatherApi/internal/dto"
	"weatherApi/internal/provider"
)

func StartSubscriptionWorker(
	ctx context.Context,
	subscriber broker.EventSubscriber,
	smtpClient provider.SMTPClientInterface,
) error {
	err := subscriber.Subscribe(ctx, broker.SendSubscriptionWeatherData, func(data []byte) error {
		var task dto.WeatherSubData
		if err := json.Unmarshal(data, &task); err != nil {
			log.Println("Failed to decode task:", err)
			return err
		}

		log.Printf("Sending subscription message to %s", task)
		return smtpClient.SendSubscriptionWeatherData(&task)
	})
	return err
}

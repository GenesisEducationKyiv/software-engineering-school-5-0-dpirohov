package worker

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"weatherApi/internal/broker"
	"weatherApi/internal/dto"
	"weatherApi/internal/provider"
)

const maxConcurrentJobs = 5

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

		var wg sync.WaitGroup
		semaphore := make(chan struct{}, maxConcurrentJobs)

		for _, user := range task.Users {
			user := user
			wg.Add(1)
			semaphore <- struct{}{}
			go func() {
				defer wg.Done()
				defer func() { <-semaphore }()

				log.Printf("Sending subscription message to user %s", user.Email)
				if err := smtpClient.SendSubscriptionWeatherData(&task.Weather, &user); err != nil {
					log.Printf("Failed to send email to %s: %v", user.Email, err)
				}
			}()
		}

		wg.Wait()
		return nil
	})
	return err
}

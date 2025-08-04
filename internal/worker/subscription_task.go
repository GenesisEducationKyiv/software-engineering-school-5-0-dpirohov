package worker

import (
	"context"
	"encoding/json"
	"sync"
	"weatherApi/internal/broker"
	"weatherApi/internal/dto"
	"weatherApi/internal/logger"
	"weatherApi/internal/provider"
)

const maxConcurrentJobs = 5

func StartSubscriptionWorker(
	log *logger.Logger,
	ctx context.Context,
	subscriber broker.EventSubscriber,
	smtpClient provider.SMTPClientInterface,
) error {
	err := subscriber.Subscribe(ctx, broker.SendSubscriptionWeatherData, func(ctx context.Context, data []byte) error {
		log := log.FromContext(ctx)

		var task dto.WeatherSubData
		if err := json.Unmarshal(data, &task); err != nil {
			log.Error().Err(err).Msg("Failed to decode task")
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
				log.Info().Msgf("Sending weather message to user %s", user.Email)
				if err := smtpClient.SendSubscriptionWeatherData(&task.Weather, &user); err != nil {
					log.Error().Err(err).Msgf("Failed to send weather email to %s", user.Email)
				}
			}()
		}

		wg.Wait()
		return nil
	})
	return err
}

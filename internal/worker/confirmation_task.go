package worker

import (
	"context"
	"encoding/json"
	"weatherApi/internal/broker"
	"weatherApi/internal/dto"
	"weatherApi/internal/logger"
	"weatherApi/internal/provider"
)

func StartConfirmationWorker(
	ctx context.Context,
	subscriber broker.EventSubscriber,
	smtpClient provider.SMTPClientInterface,
) error {
	err := subscriber.Subscribe(ctx, broker.SubscriptionConfirmationTasks, func(ctx context.Context, data []byte) error {
		log := logger.FromContext(ctx)
		var task dto.ConfirmationEmailTask
		if err := json.Unmarshal(data, &task); err != nil {
			log.Error().Err(err).Msg("Failed to decode task")
			return err
		}
		log.Info().Msgf("Sending subscription confirmation letter to %s for city %s", task.Email, task.City)
		return smtpClient.SendConfirmationToken(task.Email, task.Token, task.City)
	})
	return err
}

package worker

import (
	"encoding/json"
	"log"
	"weatherApi/internal/broker"
	"weatherApi/internal/dto"
	"weatherApi/internal/provider"
)

func StartConfirmationWorker(bus broker.EventBusInerface, smtpClient provider.SMTPClientInterface) error {
	err := bus.Subscribe(broker.SubscriptionConfirmationTasks, func(data []byte) error {
		var task dto.ConfirmationEmailTask
		if err := json.Unmarshal(data, &task); err != nil {
			log.Println("Failed to decode task:", err)
			return err
		}

		log.Printf("Sending confirmation to %s for city %s", task.Email, task.City)
		return smtpClient.SendConfirmationToken(task.Email, task.Token, task.City)
	})
	if err == nil {
		log.Printf("Successfully subscribed to topic %s", broker.SubscriptionConfirmationTasks)
	}
	return err
}

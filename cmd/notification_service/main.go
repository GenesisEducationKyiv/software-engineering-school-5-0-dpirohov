package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"weatherApi/internal/broker"
	"weatherApi/internal/config"
	"weatherApi/internal/provider"
	"weatherApi/internal/service/notification"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	cfg := config.NewNotificationServiceConfig()
	smtpClient := provider.NewSMTPClient(cfg.SmtpHost, cfg.SmtpPort, cfg.SmtpLogin, cfg.SmtpPassword, cfg.AppURL)
	publisher, err := broker.NewRabbitMQPublisher(cfg.BrokerURL)
	if err != nil {
		log.Fatalf("Publisher error: %v", err)
	}
	subscriber, err := broker.NewRabbitMQSubscriber(cfg.BrokerURL, cfg.BrokerMaxRetries, publisher)
	if err != nil {
		log.Fatalf("Subscriber error: %v", err)
	}

	err = notification.Run(ctx, notification.Service{
		Config:     cfg,
		SMTPClient: smtpClient,
		Publisher:  publisher,
		Subscriber: subscriber,
		SignalChan: sigChan,
	})
	if err != nil {
		log.Fatalf("App stopped with error: %v", err)
	}
}

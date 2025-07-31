package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"weatherApi/internal/broker"
	"weatherApi/internal/config"
	"weatherApi/internal/logger"
	"weatherApi/internal/provider"
	"weatherApi/internal/service/notification"

	"github.com/rs/zerolog"
)

func main() {
	logger.Init("notification-service", zerolog.InfoLevel)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	cfg := config.NewNotificationServiceConfig()
	smtpClient := provider.NewSMTPClient(cfg.SmtpHost, cfg.SmtpPort, cfg.SmtpLogin, cfg.SmtpPassword, cfg.AppURL)
	publisher, err := broker.NewRabbitMQPublisher(cfg.BrokerURL)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Publisher error")
	}
	subscriber, err := broker.NewRabbitMQSubscriber(cfg.BrokerURL, cfg.BrokerMaxRetries, publisher)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Subscriber error")
	}

	err = notification.Run(ctx, notification.Service{
		Config:     cfg,
		SMTPClient: smtpClient,
		Publisher:  publisher,
		Subscriber: subscriber,
		SignalChan: sigChan,
	})
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("App stopped with error")
	}
}

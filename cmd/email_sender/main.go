package email_sender

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"weatherApi/internal/broker"
	"weatherApi/internal/config"
	"weatherApi/internal/provider"
	"weatherApi/internal/worker"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go handleShutdown(cancel)

	cfg := config.LoadConfig()

	smtpClient := provider.NewSMTPClient(cfg.SmtpHost, cfg.SmtpPort, cfg.SmtpLogin, cfg.SmtpPassword, cfg.AppURL)

	publisher, err := broker.NewRabbitMQPublisher(cfg.BrokerURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ for publisher: %v", err)
	}

	subscriber, err := broker.NewRabbitMQSubscriber(cfg.BrokerURL, cfg.BrokerMaxRetries, publisher)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ for subscriber: %v", err)
	}

	log.Println("Starting email workers...")

	go func() {
		if err := worker.StartConfirmationWorker(ctx, subscriber, smtpClient); err != nil {
			log.Fatalf("ConfirmationWorker stopped with error: %v", err)
		}
	}()

	go func() {
		if err := worker.StartSubscriptionWorker(ctx, subscriber, smtpClient); err != nil {
			log.Fatalf("SubscriptionWorker stopped with error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Email service shutdown complete.")
}


func handleShutdown(cancel context.CancelFunc) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Shutdown signal received")
	cancel()
}

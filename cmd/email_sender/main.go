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
		log.Fatalf("Failed to connect to RabbitMQ for RabbitMQ subscriber: %v", err)
	}


	log.Println("Starting confirmation email worker...")

	if err := worker.StartConfirmationWorker(ctx, subscriber, smtpClient); err != nil {
		log.Fatalf("Worker stopped with error: %v", err)
	}

	log.Println("Worker shut down cleanly.")
}

func handleShutdown(cancel context.CancelFunc) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Shutdown signal received")
	cancel()
}

package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"weatherApi/internal/broker"
	"weatherApi/internal/config"
	"weatherApi/internal/provider"
	"weatherApi/internal/worker"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	cfg := config.NewNotificationServiceConfig()

	smtpClient := provider.NewSMTPClient(cfg.SmtpHost, cfg.SmtpPort, cfg.SmtpLogin, cfg.SmtpPassword, cfg.AppURL)

	publisher, err := broker.NewRabbitMQPublisher(cfg.BrokerURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ for publisher: %v", err)
	}

	subscriber, err := broker.NewRabbitMQSubscriber(cfg.BrokerURL, cfg.BrokerMaxRetries, publisher)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ for subscriber: %v", err)
	}

	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", cfg.Port),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/health" {
				w.WriteHeader(http.StatusOK)
				if _, err := w.Write([]byte("OK")); err != nil {
					return
				}
				return
			}
			http.NotFound(w, r)
		}),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	go func() {
		log.Printf("Starting health check server on :%d", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Health server error: %v", err)
		}
	}()

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

	<-sigChan
	log.Println("Shutdown signal received")
	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Health server shutdown error: %v", err)
	} else {
		log.Println("Health server shutdown gracefully")
	}

	log.Println("Email service shutdown complete.")
}

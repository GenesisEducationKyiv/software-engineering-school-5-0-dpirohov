package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
	"weatherApi/internal/broker"
	"weatherApi/internal/config"
	"weatherApi/internal/provider"
	"weatherApi/internal/worker"

	"weatherApi/internal/server"
)

func gracefulShutdown(ctx context.Context, apiServer *http.Server, rabbitMq broker.EventBusInerface, done chan bool) {
	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")

	ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := apiServer.Shutdown(ctxTimeout); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	if err := rabbitMq.Close(); err != nil {
		log.Printf("Error while closing RabbitMQ: %v", err)
	} else {
		log.Println("RabbitMQ connection closed")
	}

	log.Println("Server exiting")
	done <- true
}

func main() {
	cfg := config.LoadConfig()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	smtpClient := provider.NewSMTPClient(cfg.SmtpHost, cfg.SmtpPort, cfg.SmtpLogin, cfg.SmtpPassword, cfg.AppURL)

	rabbitMq, err := broker.NewRabbitMQBus(cfg.BrokerURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	if err := worker.StartConfirmationWorker(rabbitMq, smtpClient); err != nil {
		log.Fatalf("Failed to start confirmation worker: %v", err)
	}

	server := server.NewServer(cfg, rabbitMq)

	done := make(chan bool, 1)

	go gracefulShutdown(ctx, server, rabbitMq, done)

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("HTTP server error: %v", err)
	}

	<-done
	log.Println("Graceful shutdown complete.")
}

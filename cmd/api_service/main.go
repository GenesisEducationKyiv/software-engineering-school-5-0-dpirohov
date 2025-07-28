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
	"weatherApi/internal/scheduler"

	"weatherApi/internal/server"
)

func main() {
	cfg := config.NewApiServiceConfig()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	publisher, err := broker.NewRabbitMQPublisher(cfg.BrokerURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ for publisher: %v", err)
	}

	httpServer := server.NewServer(cfg, publisher)

	schedulerService, err := scheduler.NewService(
		httpServer.SubscriptionService.SubscriptionRepo,
		publisher,
		httpServer.WeatherService,
		ctx,
	)
	if err != nil {
		log.Fatalf("failed to init scheduler: %v", err)
	}
	if err := schedulerService.Start(ctx); err != nil {
		log.Fatalf("failed to start scheduler: %v", err)
	}

	done := make(chan bool, 1)

	go gracefulShutdown(ctx, httpServer, schedulerService, done, publisher)

	if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("HTTP server error: %v", err)
	}

	<-done
	log.Println("Graceful shutdown complete.")
}

func gracefulShutdown(ctx context.Context, httpServer *server.Server, schedulerService *scheduler.Service, done chan bool, publishers ...broker.EventPublisher) {
	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")
	ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := schedulerService.Stop(ctxTimeout); err != nil {
		log.Printf("Scheduler forced to shutdown with error: %v", err)
	}

	if err := httpServer.Shutdown(ctxTimeout); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	for i, pub := range publishers {
		if err := pub.Close(); err != nil {
			log.Printf("Error while closing RabbitMQ publisher #%d: %v", i, err)
		}
	}

	log.Println("Server exiting")
	done <- true
}

package main

import (
	"context"
	"errors"
	"net/http"
	"os/signal"
	"syscall"
	"time"
	"weatherApi/internal/broker"
	"weatherApi/internal/config"
	"weatherApi/internal/logger"
	"weatherApi/internal/scheduler"

	"github.com/rs/zerolog"

	"weatherApi/internal/server"
)

func main() {
	logger.Init("api-service", zerolog.InfoLevel)

	cfg := config.NewApiServiceConfig()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	publisher, err := broker.NewRabbitMQPublisher(cfg.BrokerURL)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to connect to RabbitMQ for publisher")
	}

	httpServer := server.NewServer(cfg, publisher)

	schedulerService, err := scheduler.NewService(
		httpServer.SubscriptionService.SubscriptionRepo,
		publisher,
		httpServer.WeatherService,
		ctx,
	)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to init scheduler")
	}
	if err := schedulerService.Start(); err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to start scheduler")
	}

	done := make(chan bool, 1)

	go gracefulShutdown(ctx, httpServer, schedulerService, done, publisher)

	if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Log.Fatal().Err(err).Msg("HTTP server error")
	}

	<-done
	logger.Log.Info().Msg("Graceful shutdown complete.")
}

func gracefulShutdown(ctx context.Context, httpServer *server.Server, schedulerService *scheduler.Service, done chan bool, publishers ...broker.EventPublisher) {
	<-ctx.Done()

	logger.Log.Warn().Msg("Shutting down gracefully, press Ctrl+C again to force")
	ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := schedulerService.Stop(); err != nil {
		logger.Log.Error().Err(err).Msg("Scheduler forced to shutdown with error")
	}

	if err := httpServer.Shutdown(ctxTimeout); err != nil {
		logger.Log.Error().Err(err).Msg("Server forced to shutdown with error")
	}

	for i, pub := range publishers {
		if err := pub.Close(); err != nil {
			logger.Log.Error().Err(err).Int("index", i).Msg("Error while closing RabbitMQ publisher")
		}
	}

	logger.Log.Info().Msg("Server exiting")
	done <- true
}

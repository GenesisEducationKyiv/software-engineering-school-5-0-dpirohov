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
	log := logger.NewLogger("api-service", zerolog.InfoLevel)

	cfg := config.NewApiServiceConfig(log.Base())

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	publisher, err := broker.NewRabbitMQPublisher(cfg.BrokerURL, log)
	if err != nil {
		log.Base().Fatal().Err(err).Msg("Failed to connect to RabbitMQ for publisher")
	}

	httpServer := server.NewServer(log, cfg, publisher)

	schedulerService, err := scheduler.NewService(
		log,
		httpServer.SubscriptionService.SubscriptionRepo,
		publisher,
		httpServer.WeatherService,
		ctx,
	)
	if err != nil {
		log.Base().Fatal().Err(err).Msg("Failed to init scheduler")
	}
	if err := schedulerService.Start(); err != nil {
		log.Base().Fatal().Err(err).Msg("Failed to start scheduler")
	}

	done := make(chan bool, 1)

	go gracefulShutdown(log, ctx, httpServer, schedulerService, done, publisher)

	if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Base().Fatal().Err(err).Msg("HTTP server error")
	}

	<-done
	log.Base().Info().Msg("Graceful shutdown complete.")
}

func gracefulShutdown(log *logger.Logger, ctx context.Context, httpServer *server.Server, schedulerService *scheduler.Service, done chan bool, publishers ...broker.EventPublisher) {
	<-ctx.Done()

	log.Base().Warn().Msg("Shutting down gracefully, press Ctrl+C again to force")
	ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := schedulerService.Stop(); err != nil {
		log.Base().Error().Err(err).Msg("Scheduler forced to shutdown with error")
	}

	if err := httpServer.Shutdown(ctxTimeout); err != nil {
		log.Base().Error().Err(err).Msg("Server forced to shutdown with error")
	}

	for i, pub := range publishers {
		if err := pub.Close(); err != nil {
			log.Base().Error().Err(err).Int("index", i).Msg("Error while closing RabbitMQ publisher")
		}
	}

	log.Base().Info().Msg("Server exiting")
	done <- true
}

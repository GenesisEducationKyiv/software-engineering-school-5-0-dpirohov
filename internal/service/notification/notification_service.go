package notification

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"
	"weatherApi/internal/broker"
	"weatherApi/internal/config"
	"weatherApi/internal/logger"
	"weatherApi/internal/provider"
	"weatherApi/internal/worker"
)

type Service struct {
	Log        *logger.Logger
	Config     *config.NotificationServiceConfig
	SMTPClient provider.SMTPClientInterface
	Publisher  broker.EventPublisher
	Subscriber broker.EventSubscriber
	SignalChan <-chan os.Signal
}

func Run(ctx context.Context, service Service) error {
	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", service.Config.Port),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/health" {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("OK"))
				return
			}
			http.NotFound(w, r)
		}),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}
	log := service.Log.Base()
	go func() {
		log.Info().Msgf("Starting health check server on :%d", service.Config.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("Health server error")
		}
	}()

	go func() {
		if err := worker.StartConfirmationWorker(service.Log, ctx, service.Subscriber, service.SMTPClient); err != nil {
			log.Fatal().Err(err).Msg("ConfirmationWorker error")
		}
	}()
	go func() {
		if err := worker.StartSubscriptionWorker(service.Log, ctx, service.Subscriber, service.SMTPClient); err != nil {
			log.Fatal().Err(err).Msg("SubscriptionWorker error")
		}
	}()

	select {
	case <-ctx.Done():
	case <-service.SignalChan:
	}

	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("Failed to shutdown notification service")
		return err
	}
	return nil
}

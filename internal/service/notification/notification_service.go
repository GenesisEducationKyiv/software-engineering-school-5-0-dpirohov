package notification

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"weatherApi/internal/broker"
	"weatherApi/internal/config"
	"weatherApi/internal/provider"
	"weatherApi/internal/worker"
)

type Service struct {
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

	go func() {
		log.Printf("Starting health check server on :%d", service.Config.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Health server error: %v", err)
		}
	}()

	go func() {
		if err := worker.StartConfirmationWorker(ctx, service.Subscriber, service.SMTPClient); err != nil {
			log.Fatalf("ConfirmationWorker error: %v", err)
		}
	}()
	go func() {
		if err := worker.StartSubscriptionWorker(ctx, service.Subscriber, service.SMTPClient); err != nil {
			log.Fatalf("SubscriptionWorker error: %v", err)
		}
	}()

	select {
	case <-ctx.Done():
	case <-service.SignalChan:
	}

	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Panicln("Failed to shutdown notification service")
		return err
	}
	return nil
}

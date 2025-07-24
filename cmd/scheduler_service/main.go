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

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"weatherApi/internal/broker"
	"weatherApi/internal/config"
	"weatherApi/internal/repository/subscription"
	schedulersvc "weatherApi/internal/scheduler"
)

//nolint:cyclop
func main() {
	cfg := config.NewSchedulerConfig()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		log.Println("Shutdown signal received")
		cancel()
	}()

	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	subRepo := subscription.NewSubscriptionRepository(db)

	publisher, err := broker.NewRabbitMQPublisher(cfg.BrokerURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	scheduler, err := schedulersvc.NewService(subRepo, publisher, cfg)
	if err != nil {
		log.Fatalf("Cannot create scheduler service: %v", err)
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
		log.Printf("Starting healthcheck server on port %d...", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Healthcheck server error: %v", err)
		}
	}()

	go func() {
		if err := scheduler.Start(ctx); err != nil {
			log.Fatalf("Failed to start scheduler: %v", err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Error shutting down healthcheck server: %v", err)
	}

	if err := scheduler.Stop(shutdownCtx); err != nil {
		log.Fatalf("Graceful shutdown of scheduler failed: %v", err)
	}

	log.Println("Scheduler shut down gracefully")
}

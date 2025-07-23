package main

import (
	"context"
	"log"
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

func main() {
	cfg := config.LoadConfig()
	ctx := context.Background()

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

	if err := scheduler.Start(ctx); err != nil {
		log.Fatalf("Failed to start scheduler: %v", err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := scheduler.Stop(shutdownCtx); err != nil {
		log.Fatalf("Graceful shutdown failed: %v", err)
	}

	log.Println("Scheduler shut down gracefully")
}

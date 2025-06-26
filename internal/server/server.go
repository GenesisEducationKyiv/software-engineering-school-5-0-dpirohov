package server

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"weatherApi/internal/broker"
	"weatherApi/internal/config"
	"weatherApi/internal/provider"

	repoSubscription "weatherApi/internal/repository/subscription"
	repoUser "weatherApi/internal/repository/user"
	serviceHealthcheck "weatherApi/internal/service/healthcheck"
	serviceSubscription "weatherApi/internal/service/subscription"
	serviceWeather "weatherApi/internal/service/weather"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Server struct {
	config              *config.Config
	WeatherService      *serviceWeather.Service
	SubscriptionService *serviceSubscription.SubscriptionService
	HealthCheckService  serviceHealthcheck.HealthCheckService
}

func NewServer(cfg *config.Config, broker broker.EventPublisher) *http.Server {

	gormDB, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		log.Fatalf("Failed to get sql.DB: %v", err)
	}

	userRepo := repoUser.NewUserRepository(gormDB)
	subscriptionRepo := repoSubscription.NewSubscriptionRepository(gormDB)

	weatherService := serviceWeather.NewWeatherService(
		provider.NewOpenWeatherApiProvider(cfg.OpenWeatherAPIkey, "http://api.openweathermap.org/data/2.5/weather"),
		provider.NewWeatherApiProvider(cfg.WeatherApiAPIkey, "http://api.weatherapi.com/v1/current.json"),
	)
	subscriptionService := serviceSubscription.NewSubscriptionService(
		subscriptionRepo,
		userRepo,
		broker,
		cfg.TokenLifetimeMinutes,
	)
	healthcheckService := serviceHealthcheck.New(sqlDB)

	NewServer := &Server{
		config:              cfg,
		WeatherService:      weatherService,
		SubscriptionService: subscriptionService,
		HealthCheckService:  healthcheckService,
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.config.Port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}

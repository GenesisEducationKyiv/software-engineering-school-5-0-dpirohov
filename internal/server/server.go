package server

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"weatherApi/internal/broker"
	"weatherApi/internal/config"
	"weatherApi/internal/metrics"
	"weatherApi/internal/provider"
	"weatherApi/internal/repository/weather"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/redis/go-redis/v9"

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
	config              *config.ApiServiceConfig
	WeatherService      *serviceWeather.Service
	SubscriptionService *serviceSubscription.SubscriptionService
	HealthCheckService  serviceHealthcheck.HealthCheckService
}

func NewServer(cfg *config.ApiServiceConfig, broker broker.EventPublisher) *http.Server {

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

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisURL,
		Password: cfg.RedisPassword,
		DB:       0,
	})

	userRepo := repoUser.NewUserRepository(gormDB)
	subscriptionRepo := repoSubscription.NewSubscriptionRepository(gormDB)
	cacheMetrics := metrics.NewCacheMetrics()
	cacheMetrics.Register(prometheus.DefaultRegisterer)
	cacheRepo := weather.NewWeatherRepository(&weather.RepositoryOptions{
		Client:       rdb,
		CacheTTL:     cfg.CacheTTL,
		LockTTL:      cfg.LockTTL,
		LockRetryDur: cfg.LockRetryDur,
		LockMaxWait:  cfg.LockMaxWait,
		Metrics:      cacheMetrics,
	})

	weatherService := serviceWeather.NewWeatherService(
		cacheRepo,
		provider.NewOpenWeatherApiProvider(cfg.OpenWeatherAPIkey, cfg.OpenWeatherAPIEndpoint),
		provider.NewWeatherApiProvider(cfg.WeatherApiAPIkey, cfg.WeatherApiAPIEndpoint),
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

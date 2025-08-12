package server

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"weatherApi/internal/broker"
	"weatherApi/internal/config"
	"weatherApi/internal/logger"
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
	gormLogger "gorm.io/gorm/logger"
)

type Server struct {
	log                 *logger.Logger
	config              *config.ApiServiceConfig
	WeatherService      *serviceWeather.Service
	SubscriptionService *serviceSubscription.SubscriptionService
	HealthCheckService  serviceHealthcheck.HealthCheckService
	httpServer          *http.Server
}

func NewServer(log *logger.Logger, cfg *config.ApiServiceConfig, broker broker.EventPublisher) *Server {

	gormDB, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Silent),
	})
	if err != nil {
		log.Base().Fatal().Err(err).Msg("Failed to connect to DB")
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		log.Base().Fatal().Err(err).Msg("Failed to get sql.DB")
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
		log,
		cacheRepo,
		provider.NewOpenWeatherApiProvider(log, cfg.OpenWeatherAPIkey, cfg.OpenWeatherAPIEndpoint),
		provider.NewWeatherApiProvider(log, cfg.WeatherApiAPIkey, cfg.WeatherApiAPIEndpoint),
	)
	subscriptionService := serviceSubscription.NewSubscriptionService(
		log,
		subscriptionRepo,
		userRepo,
		broker,
		cfg.TokenLifetimeMinutes,
	)
	healthcheckService := serviceHealthcheck.New(log, sqlDB)

	server := &Server{
		log:                 log,
		config:              cfg,
		WeatherService:      weatherService,
		SubscriptionService: subscriptionService,
		HealthCheckService:  healthcheckService,
	}

	server.httpServer = &http.Server{
		Addr:         fmt.Sprintf(":%d", server.config.Port),
		Handler:      server.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}

func (s *Server) ListenAndServe() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

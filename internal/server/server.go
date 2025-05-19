package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"weatherApi/internal/provider"
	repo_subscription "weatherApi/internal/repository/subscription"
	repo_user "weatherApi/internal/repository/user"

	s_healthcheck "weatherApi/internal/service/healthcheck"
	s_subscription "weatherApi/internal/service/subscription"
	s_weather "weatherApi/internal/service/weather"

	_ "github.com/joho/godotenv/autoload"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Server struct {
	port                int
	WeatherService      *s_weather.WeatherService
	SubscriptionService *s_subscription.SubscriptionService
	HealthCheckService  s_healthcheck.HealthCheckService
}

func NewServer() *http.Server {
	host := os.Getenv("HOST")
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	gormDB, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})

	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		log.Fatalf("Failed to get sql.DB: %v", err)
	}

	userRepo := repo_user.NewUserRepository(gormDB)
	subscriptionRepo := repo_subscription.NewSubscriptionRepository(gormDB)

	weatherService := s_weather.NewWeatherService()

	jwtLifetimeMinutes, err := strconv.Atoi(os.Getenv("TOKEN_LIFETIME_MINUTES"))

	if err != nil {
		log.Fatal("TOKEN_LIFETIME_MINUTES not provided!")
	}

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")

	if smtpHost == "" || smtpUser == "" || smtpPass == "" || err != nil {
		log.Fatal("Fail to retrieve SMTP credentials")
	}

	smtpClient := provider.NewSMTPClient(smtpHost, smtpPort, smtpUser, smtpPass, fmt.Sprintf("%s:%d", host, port))
	subscriptionService := s_subscription.NewSubscriptionService(
		subscriptionRepo,
		userRepo,
		smtpClient,
		jwtLifetimeMinutes,
	)

	healthcheckService := s_healthcheck.New(sqlDB)

	NewServer := &Server{
		port:                port,
		WeatherService:      weatherService,
		SubscriptionService: subscriptionService,
		HealthCheckService:  healthcheckService,
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}

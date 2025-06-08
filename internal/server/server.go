package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"weatherApi/internal/provider"
	repoSubscription "weatherApi/internal/repository/subscription"
	repoUser "weatherApi/internal/repository/user"
	serviceHealthcheck "weatherApi/internal/service/healthcheck"
	serviceSubscription "weatherApi/internal/service/subscription"
	serviceWeather "weatherApi/internal/service/weather"

	// Automatically loads environment variables from .env file on startup
	_ "github.com/joho/godotenv/autoload"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Server struct {
	port                int
	WeatherService      *serviceWeather.WeatherService
	SubscriptionService *serviceSubscription.SubscriptionService
	HealthCheckService  serviceHealthcheck.HealthCheckService
}

func NewServer() *http.Server {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatal("Server port is not provided!")
	}
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

	userRepo := repoUser.NewUserRepository(gormDB)
	subscriptionRepo := repoSubscription.NewSubscriptionRepository(gormDB)

	weatherService := serviceWeather.NewWeatherService()

	jwtLifetimeMinutes, err := strconv.Atoi(os.Getenv("TOKEN_LIFETIME_MINUTES"))
	if err != nil {
		log.Fatal("TOKEN_LIFETIME_MINUTES not provided!")
	}

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")
	smtpLinkUrl := os.Getenv("SMTP_LINK_URL")

	if smtpHost == "" || smtpUser == "" || smtpPass == "" || smtpLinkUrl == "" || err != nil {
		log.Fatal("Fail to retrieve SMTP credentials")
	}

	smtpClient := provider.NewSMTPClient(smtpHost, smtpPort, smtpUser, smtpPass, smtpLinkUrl)
	subscriptionService := serviceSubscription.NewSubscriptionService(
		subscriptionRepo,
		userRepo,
		smtpClient,
		jwtLifetimeMinutes,
	)

	healthcheckService := serviceHealthcheck.New(sqlDB)

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

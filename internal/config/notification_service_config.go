package config

import (
	"path/filepath"
	"weatherApi/internal/logger"

	"github.com/joho/godotenv"
)

type NotificationServiceConfig struct {
	Config
	Port             int
	AppURL           string
	BrokerURL        string
	BrokerMaxRetries int

	SmtpHost     string
	SmtpPort     int
	SmtpLogin    string
	SmtpPassword string

	RootDir string
}

func NewNotificationServiceConfig() *NotificationServiceConfig {
	rootDir := getRootDir()
	err := godotenv.Load(filepath.Join(rootDir, ".env.notification_service"))
	if err != nil {
		logger.Log.Warn().Msg("Failed to load .env file!")
	}
	return &NotificationServiceConfig{
		Port:             mustGet[int]("PORT"),
		AppURL:           mustGet[string]("APP_URL"),
		BrokerURL:        mustGet[string]("BROKER_URL"),
		BrokerMaxRetries: getWithDefault[int]("RMQ_MAX_RETRIES", 3),
		SmtpHost:         mustGet[string]("SMTP_HOST"),
		SmtpPort:         mustGet[int]("SMTP_PORT"),
		SmtpLogin:        mustGet[string]("SMTP_USER"),
		SmtpPassword:     mustGet[string]("SMTP_PASS"),
		RootDir:          rootDir,
	}
}

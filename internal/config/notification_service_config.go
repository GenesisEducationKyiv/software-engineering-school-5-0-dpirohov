package config

import (
	"path/filepath"

	"github.com/rs/zerolog"

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

func NewNotificationServiceConfig(log *zerolog.Logger) *NotificationServiceConfig {
	rootDir := getRootDir(log)
	err := godotenv.Load(filepath.Join(rootDir, ".env.notification_service"))
	if err != nil {
		log.Warn().Msg("Failed to load .env file!")
	}
	return &NotificationServiceConfig{
		Port:             mustGet[int](log, "PORT"),
		AppURL:           mustGet[string](log, "APP_URL"),
		BrokerURL:        mustGet[string](log, "BROKER_URL"),
		BrokerMaxRetries: getWithDefault[int](log, "RMQ_MAX_RETRIES", 3),
		SmtpHost:         mustGet[string](log, "SMTP_HOST"),
		SmtpPort:         mustGet[int](log, "SMTP_PORT"),
		SmtpLogin:        mustGet[string](log, "SMTP_USER"),
		SmtpPassword:     mustGet[string](log, "SMTP_PASS"),
		RootDir:          rootDir,
	}
}

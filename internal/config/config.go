package config

import (
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Host string
	Port int

	AppURL           string
	DatabaseURL      string
	BrokerURL        string
	BrokerMaxRetries int

	OpenWeatherAPIkey string
	WeatherApiAPIkey  string

	TokenLifetimeMinutes int

	SmtpHost     string
	SmtpPort     int
	SmtpLogin    string
	SmtpPassword string

	RootDir string
}

func LoadConfig() *Config {
	rootDir := getRootDir()
	err := godotenv.Load(filepath.Join(rootDir, ".env"))
	if err != nil {
		log.Printf("Failed to load .env file! Err: %v", err)
	}
	return &Config{
		Host:                 mustGet[string]("HOST"),
		Port:                 mustGet[int]("PORT"),
		AppURL:               mustGet[string]("APP_URL"),
		DatabaseURL:          mustGet[string]("DB_URL"),
		BrokerURL:            mustGet[string]("BROKER_URL"),
		BrokerMaxRetries:     getWithDefault[int]("RMQ_MAX_RETRIES", 3),
		OpenWeatherAPIkey:    mustGet[string]("OPENWEATHER_API_KEY"),
		WeatherApiAPIkey:     mustGet[string]("WEATHER_API_API_KEY"),
		TokenLifetimeMinutes: getWithDefault[int]("TOKEN_LIFETIME_MINUTES", 15),
		SmtpHost:             mustGet[string]("SMTP_HOST"),
		SmtpPort:             mustGet[int]("SMTP_PORT"),
		SmtpLogin:            mustGet[string]("SMTP_USER"),
		SmtpPassword:         mustGet[string]("SMTP_PASS"),
		RootDir:              rootDir,
	}
}

func mustGet[T any](key string) T {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("missing required environment variable: %s", key)
	}
	return castEnvValue[T](val, key)
}

func getWithDefault[T any](key string, defaultVal T) T {
	val := os.Getenv(key)
	if val == "" {
		log.Printf("missing optional environment variable: %s, using default value: %v", key, defaultVal)
		return defaultVal
	}
	return castEnvValue[T](val, key)
}

func castEnvValue[T any](val string, key string) T {
	var zero T

	switch any(zero).(type) {
	case string:
		return any(val).(T)
	case int:
		intVal, err := strconv.Atoi(val)
		if err != nil {
			log.Fatalf("invalid int value for %s: %v", key, err)
		}
		return any(intVal).(T)
	default:
		log.Fatalf("unsupported type for env variable: %T", zero)
	}

	return zero // unreachable, но нужен для компиляции
}

func getRootDir() string {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("cannot get working directory: %v", err)
	}
	if os.Getenv("ENV") != "DOCKER" {
		return filepath.Join(currentDir, "../../")
	}
	return currentDir
}

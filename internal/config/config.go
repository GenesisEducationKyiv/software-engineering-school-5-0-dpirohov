package config

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Host string
	Port int

	AppURL           string
	DatabaseURL      string
	BrokerURL        string
	BrokerMaxRetries int

	OpenWeatherAPIEndpoint string
	OpenWeatherAPIkey      string
	WeatherApiAPIEndpoint  string
	WeatherApiAPIkey       string

	TokenLifetimeMinutes int

	SmtpHost     string
	SmtpPort     int
	SmtpLogin    string
	SmtpPassword string

	RootDir string

	RedisUlr      string
	RedisPassword string
	CacheTTL      time.Duration
	LockTTL       time.Duration
	LockRetryDur  time.Duration
	LockMaxWait   time.Duration
}

func LoadConfig() *Config {
	rootDir := getRootDir()
	err := godotenv.Load(filepath.Join(rootDir, ".env"))
	if err != nil {
		log.Printf("Failed to load .env file! Err: %v", err)
	}
	return &Config{
		Host:                   mustGet[string]("HOST"),
		Port:                   mustGet[int]("PORT"),
		AppURL:                 mustGet[string]("APP_URL"),
		DatabaseURL:            mustGet[string]("DB_URL"),
		BrokerURL:              mustGet[string]("BROKER_URL"),
		BrokerMaxRetries:       getWithDefault[int]("RMQ_MAX_RETRIES", 3),
		OpenWeatherAPIEndpoint: mustGet[string]("OPENWEATHER_API_ENDPOINT"),
		OpenWeatherAPIkey:      mustGet[string]("OPENWEATHER_API_KEY"),
		WeatherApiAPIEndpoint:  mustGet[string]("WEATHER_API_API_ENDPOINT"),
		WeatherApiAPIkey:       mustGet[string]("WEATHER_API_API_KEY"),
		TokenLifetimeMinutes:   getWithDefault[int]("TOKEN_LIFETIME_MINUTES", 15),
		SmtpHost:               mustGet[string]("SMTP_HOST"),
		SmtpPort:               mustGet[int]("SMTP_PORT"),
		SmtpLogin:              mustGet[string]("SMTP_USER"),
		SmtpPassword:           mustGet[string]("SMTP_PASS"),
		RootDir:                rootDir,
		RedisUlr:               mustGet[string]("REDIS_URL"),
		RedisPassword:          mustGet[string]("REDIS_PWD"),
		CacheTTL:               getWithDefault[time.Duration]("CACHE_TTL", 5*time.Minute),
		LockTTL:                getWithDefault[time.Duration]("LOCK_TTL", 3*time.Second),
		LockRetryDur:           getWithDefault[time.Duration]("LOCK_RETRY_DUR", 100*time.Millisecond),
		LockMaxWait:            getWithDefault[time.Duration]("LOCK_MAX_WAIT", 3*time.Second),
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
	case time.Duration:
		dur, err := time.ParseDuration(val)
		if err != nil {
			log.Fatalf("invalid duration value for %s: %v", key, err)
		}
		return any(dur).(T)
	default:
		log.Fatalf("unsupported type for env variable: %T", zero)
	}

	return zero
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

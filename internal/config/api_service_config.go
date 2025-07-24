package config

import (
	"log"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
)

type ApiServiceConfig struct {
	Config
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
	TokenLifetimeMinutes   int

	RootDir string

	RedisURL      string
	RedisPassword string
	CacheTTL      time.Duration
	LockTTL       time.Duration
	LockRetryDur  time.Duration
	LockMaxWait   time.Duration
}

func NewApiServiceConfig() *ApiServiceConfig {
	rootDir := getRootDir()
	err := godotenv.Load(filepath.Join(rootDir, ".env.api_service"))
	if err != nil {
		log.Printf("Failed to load .env file! Err: %v", err)
	}
	return &ApiServiceConfig{
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
		RootDir:                rootDir,
		RedisURL:               mustGet[string]("REDIS_URL"),
		RedisPassword:          mustGet[string]("REDIS_PWD"),
		CacheTTL:               getWithDefault[time.Duration]("CACHE_TTL", 5*time.Minute),
		LockTTL:                getWithDefault[time.Duration]("LOCK_TTL", 3*time.Second),
		LockRetryDur:           getWithDefault[time.Duration]("LOCK_RETRY_DUR", 100*time.Millisecond),
		LockMaxWait:            getWithDefault[time.Duration]("LOCK_MAX_WAIT", 3*time.Second),
	}
}

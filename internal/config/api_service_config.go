package config

import (
	"path/filepath"
	"time"

	"github.com/rs/zerolog"

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

func NewApiServiceConfig(log *zerolog.Logger) *ApiServiceConfig {
	rootDir := getRootDir(log)
	err := godotenv.Load(filepath.Join(rootDir, ".env.api_service"))
	if err != nil {
		log.Error().Err(err).Msg("Failed to load .env file!")
	}
	return &ApiServiceConfig{
		Host:                   mustGet[string](log, "HOST"),
		Port:                   mustGet[int](log, "PORT"),
		AppURL:                 mustGet[string](log, "APP_URL"),
		DatabaseURL:            mustGet[string](log, "DB_URL"),
		BrokerURL:              mustGet[string](log, "BROKER_URL"),
		BrokerMaxRetries:       getWithDefault[int](log, "RMQ_MAX_RETRIES", 3),
		OpenWeatherAPIEndpoint: mustGet[string](log, "OPENWEATHER_API_ENDPOINT"),
		OpenWeatherAPIkey:      mustGet[string](log, "OPENWEATHER_API_KEY"),
		WeatherApiAPIEndpoint:  mustGet[string](log, "WEATHER_API_API_ENDPOINT"),
		WeatherApiAPIkey:       mustGet[string](log, "WEATHER_API_API_KEY"),
		TokenLifetimeMinutes:   getWithDefault[int](log, "TOKEN_LIFETIME_MINUTES", 15),
		RootDir:                rootDir,
		RedisURL:               mustGet[string](log, "REDIS_URL"),
		RedisPassword:          mustGet[string](log, "REDIS_PWD"),
		CacheTTL:               getWithDefault[time.Duration](log, "CACHE_TTL", 5*time.Minute),
		LockTTL:                getWithDefault[time.Duration](log, "LOCK_TTL", 3*time.Second),
		LockRetryDur:           getWithDefault[time.Duration](log, "LOCK_RETRY_DUR", 100*time.Millisecond),
		LockMaxWait:            getWithDefault[time.Duration](log, "LOCK_MAX_WAIT", 3*time.Second),
	}
}

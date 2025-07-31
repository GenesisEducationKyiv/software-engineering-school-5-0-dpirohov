package config

import (
	"os"
	"path/filepath"
	"strconv"
	"time"
	"weatherApi/internal/logger"
)

type Config struct{}

func mustGet[T any](key string) T {
	val := os.Getenv(key)
	if val == "" {
		logger.Log.Fatal().Msgf("Missing required environment variable: %s", key)
	}
	return castEnvValue[T](val, key)
}

func getWithDefault[T any](key string, defaultVal T) T {
	val := os.Getenv(key)
	if val == "" {
		logger.Log.Warn().Msgf("Missing optional environment variable: %s, using default value: %v", key, defaultVal)
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
			logger.Log.Fatal().Err(err).Msgf("Invalid int value for %s", key)
		}
		return any(intVal).(T)
	case time.Duration:
		dur, err := time.ParseDuration(val)
		if err != nil {
			logger.Log.Fatal().Err(err).Msgf("Invalid duration value for %s", key)
		}
		return any(dur).(T)
	default:
		logger.Log.Fatal().Msgf("unsupported type for env variable: %T", zero)
	}

	return zero
}

func getRootDir() string {
	currentDir, err := os.Getwd()
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Cannot get working directory")
	}
	if os.Getenv("ENV") != "DOCKER" {
		return filepath.Join(currentDir, "../../")
	}
	return currentDir
}

package config

import (
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/rs/zerolog"
)

type Config struct{}

func mustGet[T any](log *zerolog.Logger, key string) T {
	val := os.Getenv(key)
	if val == "" {
		log.Fatal().Msgf("Missing required environment variable: %s", key)
	}
	return castEnvValue[T](log, val, key)
}

func getWithDefault[T any](log *zerolog.Logger, key string, defaultVal T) T {
	val := os.Getenv(key)
	if val == "" {
		log.Warn().Msgf("Missing optional environment variable: %s, using default value: %v", key, defaultVal)
		return defaultVal
	}
	return castEnvValue[T](log, val, key)
}

func castEnvValue[T any](log *zerolog.Logger, val string, key string) T {
	var zero T

	switch any(zero).(type) {
	case string:
		return any(val).(T)
	case int:
		intVal, err := strconv.Atoi(val)
		if err != nil {
			log.Fatal().Err(err).Msgf("Invalid int value for %s", key)
		}
		return any(intVal).(T)
	case time.Duration:
		dur, err := time.ParseDuration(val)
		if err != nil {
			log.Fatal().Err(err).Msgf("Invalid duration value for %s", key)
		}
		return any(dur).(T)
	default:
		log.Fatal().Msgf("unsupported type for env variable: %T", zero)
	}

	return zero
}

func getRootDir(log *zerolog.Logger) string {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot get working directory")
	}
	if os.Getenv("ENV") != "DOCKER" {
		return filepath.Join(currentDir, "../../")
	}
	return currentDir
}

package config

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type Config struct{}

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

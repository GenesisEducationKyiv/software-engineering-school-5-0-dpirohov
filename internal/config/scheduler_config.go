package config

import (
	"log"
	"path/filepath"

	"github.com/joho/godotenv"
)

type SchedulerConfig struct {
	Config
	Port               int
	APIServiceEndpoint string
	DatabaseURL        string
	BrokerURL          string
	BrokerMaxRetries   int
	RootDir            string
}

func NewSchedulerConfig() *SchedulerConfig {
	rootDir := getRootDir()
	err := godotenv.Load(filepath.Join(rootDir, ".env.scheduler"))
	if err != nil {
		log.Printf("Failed to load .env file! Err: %v", err)
	}
	return &SchedulerConfig{
		Port:               mustGet[int]("PORT"),
		APIServiceEndpoint: mustGet[string]("API_ENDPOINT"),
		DatabaseURL:        mustGet[string]("DB_URL"),
		BrokerURL:          mustGet[string]("BROKER_URL"),
		BrokerMaxRetries:   getWithDefault[int]("RMQ_MAX_RETRIES", 3),
		RootDir:            rootDir,
	}
}

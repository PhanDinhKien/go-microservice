package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	ServiceName string
	Environment string
}

func LoadConfig() *Config {
	// Load .env file if exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	return &Config{
		Port:        getEnv("PORT", "8080"),
		ServiceName: getEnv("SERVICE_NAME", "microservice"),
		Environment: getEnv("ENVIRONMENT", "development"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

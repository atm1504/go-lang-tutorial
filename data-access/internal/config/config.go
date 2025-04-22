package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}

	config := &Config{
		DBUser:     getEnvOrDefault("DBUSER", "root"),
		DBPassword: getEnvOrDefault("DBPASS", "11312113"),
		DBHost:     getEnvOrDefault("DBHOST", "127.0.0.1"),
		DBPort:     getEnvOrDefault("DBPORT", "3306"),
		DBName:     getEnvOrDefault("DBNAME", "recordings"),
	}

	return config, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

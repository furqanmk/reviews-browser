package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds the application configuration.
type Config struct {
	AppStoreReviewsURL string
	ServerPort         string
	ReviewsCSV         string
	AppsCSV            string
	RecencyCutoffHrs   int
	CleanupEveryHrs    int
}

// LoadEnv loads env vars from .env file
func LoadEnv() error {
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("error loading .env file")
	}
	return nil
}

// GetConfig returns the configuration loaded from environment variables.
func GetConfig() (*Config, error) {
	return &Config{
		AppStoreReviewsURL: os.Getenv("APPSTORE_REVIEW_URL"),
		ServerPort:         os.Getenv("SERVER_PORT"),
		ReviewsCSV:         os.Getenv("REVIEWS_CSV_PATH"),
		AppsCSV:            os.Getenv("APPS_CSV_PATH"),
		RecencyCutoffHrs:   getEnvAsInt("RECENCY_CUTOFF_HRS", 48),
		CleanupEveryHrs:    getEnvAsInt("CLEANUP_EVERY_HRS", 1),
	}, nil
}

func getEnvAsInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intValue
}

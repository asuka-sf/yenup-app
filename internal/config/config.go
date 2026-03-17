package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort            string
	BaseCurrency       string
	TargetCurrency     string
	APIProvider        string // "frankfurter" or "exchangerates"
	ExchangeRateAPIKey string
	ExchangeRateAPIURL string
	FrankfurterAPIURL  string
	SlackWebhookURL    string
	GCSBucketName      string
	GCSObjectName      string
}

func Load() (*Config, error) {
	// load the config from the environment variables
	_ = godotenv.Load()

	cfg := &Config{
		// Cloud Run sets PORT, but we also support APP_PORT for local dev
		AppPort:            getEnv("PORT", getEnv("APP_PORT", "8080")),
		BaseCurrency:       getEnv("BASE_CURRENCY", "CAD"),
		TargetCurrency:     getEnv("TARGET_CURRENCY", "JPY"),
		APIProvider:        getEnv("API_PROVIDER", "frankfurter"), // Default to frankfurter (free, no API key)
		ExchangeRateAPIKey: getEnv("EXCHANGE_RATE_API_KEY", ""),
		ExchangeRateAPIURL: getEnv("EXCHANGE_RATE_API_URL", ""),
		FrankfurterAPIURL:  getEnv("FRANKFURTER_API_URL", "https://api.frankfurter.app/"),
		SlackWebhookURL:    getEnv("SLACK_WEBHOOK_URL", ""),
		GCSBucketName:      getEnv("GCS_BUCKET_NAME", ""),
		GCSObjectName:      getEnv("GCS_OBJECT_NAME", ""),
	}
	return cfg, nil
}

func getEnv(key string, fallback string) string {
	// return the value of the environment variable if it exists, otherwise return an empty string
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

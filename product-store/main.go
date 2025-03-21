package main

import (
	"os"
	"product-store/pkg/xredis"

	"github.com/rs/zerolog"
)

func main() {
	// Get Redis connection details from environment variables
	redisHost := getEnv("REDIS_HOST", "localhost")
	redisPort := getEnv("REDIS_PORT", "6379")
	redisPassword := getEnv("REDIS_PASSWORD", "")

	logger := zerolog.New(os.Stdout).With().Timestamp().Caller().Logger()

	opts := xredis.ClientOpts{
		Host:     redisHost,
		Port:     redisPort,
		Password: redisPassword,
		Logger:   &logger,
	}

	// Set up Redis client using xredis package
	redisClient := xredis.NewClient(opts)

	// Create new handler with Redis client
	handler := NewHandler(&logger, redisClient)

	// Start server
	handler.ListenAndServe()

	// Start server
	handler.ListenAndServe()
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

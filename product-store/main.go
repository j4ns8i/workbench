package main

import (
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"

	"product-store/pkg/api"
	"product-store/pkg/xredis"
)

func main() {
	// Get Redis connection details from environment variables
	redisHost := getEnv("REDIS_HOST", "localhost")
	redisPort := getEnv("REDIS_PORT", "6379")
	redisPassword := getEnv("REDIS_PASSWORD", "")

	logger := zerolog.New(os.Stdout).With().Timestamp().Caller().Logger()

	// Set up Redis client using xredis package
	redisClient := newRedisClient(redisHost, redisPort, redisPassword)
	db := xredis.NewDB(redisClient, &logger)

	// Create new handler with Redis client
	handler := api.NewHandler(&logger, db)

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

func newRedisClient(host, port, password string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:        fmt.Sprintf("%s:%s", host, port),
		Password:    password,
		DB:          0,
		Protocol:    3,
		MaxRetries:  10,
		DialTimeout: 1 * time.Second,
	})
}

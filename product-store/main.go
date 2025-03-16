package main

import (
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

func main() {
	// Get Redis connection details from environment variables
	redisHost := getEnv("REDIS_HOST", "localhost")
	redisPort := getEnv("REDIS_PORT", "6379")
	redisPassword := getEnv("REDIS_PASSWORD", "")

	// Set up Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: redisPassword,
		DB:       0, // use default DB
		Protocol: 3,
	})

	// Create new handler with Redis client
	handler := NewHandler(rdb)

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

package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds all application configuration
type Config struct {
	Port            string
	MaxUploadSizeMB int64
	RateLimit       int
	APIKey          string
	ShutdownTimeout time.Duration
	RequestTimeout  time.Duration
	WorkerPoolSize  int
	LogLevel        string
}

// Load reads configuration from environment variables with sensible defaults
func Load() *Config {
	return &Config{
		Port:            getEnv("PORT", "8080"),
		MaxUploadSizeMB: getEnvAsInt64("MAX_UPLOAD_SIZE_MB", 10),
		RateLimit:       getEnvAsInt("RATE_LIMIT", 100),
		APIKey:          getEnv("API_KEY", ""),
		ShutdownTimeout: getEnvAsDuration("SHUTDOWN_TIMEOUT", 30*time.Second),
		RequestTimeout:  getEnvAsDuration("REQUEST_TIMEOUT", 60*time.Second),
		WorkerPoolSize:  getEnvAsInt("WORKER_POOL_SIZE", 10),
		LogLevel:        getEnv("LOG_LEVEL", "info"),
	}
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt retrieves an environment variable as int or returns default
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

// getEnvAsInt64 retrieves an environment variable as int64 or returns default
func getEnvAsInt64(key string, defaultValue int64) int64 {
	valueStr := os.Getenv(key)
	if value, err := strconv.ParseInt(valueStr, 10, 64); err == nil {
		return value
	}
	return defaultValue
}

// getEnvAsDuration retrieves an environment variable as duration or returns default
func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := os.Getenv(key)
	if value, err := time.ParseDuration(valueStr); err == nil {
		return value
	}
	return defaultValue
}

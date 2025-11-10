package config

import (
	"os"
	"strconv"
	"time"
)

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

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsInt64(key string, defaultValue int64) int64 {
	valueStr := os.Getenv(key)
	if value, err := strconv.ParseInt(valueStr, 10, 64); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := os.Getenv(key)
	if value, err := time.ParseDuration(valueStr); err == nil {
		return value
	}
	return defaultValue
}

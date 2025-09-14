package config

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server struct {
		Host         string
		Port         string
		ReadTimeout  time.Duration
		WriteTimeout time.Duration
		IdleTimeout  time.Duration
	}
	Database struct {
		Host         string
		Port         string
		User         string
		Password     string
		DBName       string
		SSLMode      string
		MaxOpenConns int
		MaxIdleConns int
	}
	MLService struct {
		BaseURL string
		Timeout time.Duration
		APIKey  string
	}
	Storage struct {
		UploadPath   string
		MaxFileSize  int64
		AllowedTypes []string
	}
	Queue struct {
		Enabled   bool
		Type      string
		URL       string
		QueueName string
	}
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{}

	cfg.Server.Host = getEnv("SERVER_HOST", "localhost")
	cfg.Server.Port = getEnv("SERVER_PORT", "8080")
	cfg.Server.ReadTimeout = getEnvDuration("SERVER_READ_TIMEOUT", "15s")
	cfg.Server.WriteTimeout = getEnvDuration("SERVER_WRITE_TIMEOUT", "15s")
	cfg.Server.IdleTimeout = getEnvDuration("SERVER_IDLE_TIMEOUT", "60s")

	cfg.Database.Host = getEnv("DB_HOST", "localhost")
	cfg.Database.Port = getEnv("DB_PORT", "5432")
	cfg.Database.User = getEnv("DB_USER", "postgres")
	cfg.Database.Password = getEnv("DB_PASSWORD", "")
	cfg.Database.DBName = getEnv("DB_NAME", "car_status_db")
	cfg.Database.SSLMode = getEnv("DB_SSLMODE", "disable")
	cfg.Database.MaxOpenConns = getEnvInt("DB_MAX_OPEN_CONNS", 25)
	cfg.Database.MaxIdleConns = getEnvInt("DB_MAX_IDLE_CONNS", 5)

	cfg.MLService.BaseURL = getEnv("ML_SERVICE_URL", "http://localhost:8000")
	cfg.MLService.Timeout = getEnvDuration("ML_SERVICE_TIMEOUT", "30s")
	cfg.MLService.APIKey = getEnv("ML_SERVICE_API_KEY", "")

	cfg.Storage.UploadPath = getEnv("UPLOAD_PATH", "./uploads")
	cfg.Storage.MaxFileSize = getEnvInt64("MAX_FILE_SIZE", 10485760) // 10MB
	allowedTypesStr := getEnv("ALLOWED_TYPES", "image/jpeg,image/jpg,image/png")
	cfg.Storage.AllowedTypes = strings.Split(allowedTypesStr, ",")

	cfg.Queue.Enabled = getEnvBool("QUEUE_ENABLED", false)
	cfg.Queue.Type = getEnv("QUEUE_TYPE", "db")
	cfg.Queue.URL = getEnv("QUEUE_URL", "")
	cfg.Queue.QueueName = getEnv("QUEUE_NAME", "prediction_jobs")

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvInt64(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getEnvDuration(key, defaultValue string) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	duration, _ := time.ParseDuration(defaultValue)
	return duration
}
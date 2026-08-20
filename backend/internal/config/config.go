package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	AppPort     string
	AppEnv      string
	JWTSecret   string
	JWTExpiry   time.Duration
	RefreshExpiry time.Duration

	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string

	MinIOEndpoint  string
	MinIOAccessKey string
	MinIOSecretKey string
	MinIOBucket    string
	MinIORegion    string
	MinIOUseSSL    bool

	BrevoAPIKey    string
	BrevoSender    string
	BrevoSenderEmail string

	FrontendURL   string
	SuperAdminEmail string
	SuperAdminPassword string
}

func Load() *Config {
	return &Config{
		AppPort:     getEnv("APP_PORT", "8083"),
		AppEnv:      getEnv("APP_ENV", "production"),
		JWTSecret:   getEnv("JWT_SECRET", "change-me-in-production"),
		JWTExpiry:   getDuration("JWT_EXPIRY_MIN", 15),
		RefreshExpiry: getDuration("REFRESH_EXPIRY_HOUR", 168), // 7 days

		DBHost:     getEnv("DB_HOST", "postgres"),
		DBPort:     getInt("DB_PORT", 5432),
		DBUser:     getEnv("DB_USER", "klikku"),
		DBPassword: getEnv("DB_PASSWORD", "klikku"),
		DBName:     getEnv("DB_NAME", "klikku"),

		MinIOEndpoint:  getEnv("MINIO_ENDPOINT", "minio:8089"),
		MinIOAccessKey: getEnv("MINIO_ACCESS_KEY", "minioadmin"),
		MinIOSecretKey: getEnv("MINIO_SECRET_KEY", "minioadmin"),
		MinIOBucket:    getEnv("MINIO_BUCKET", "klikku"),
		MinIORegion:    getEnv("MINIO_REGION", "us-east-1"),
		MinIOUseSSL:    getBool("MINIO_USE_SSL", false),

		BrevoAPIKey:    getEnv("BREVO_API_KEY", ""),
		BrevoSender:    getEnv("BREVO_SENDER", "Klikku"),
		BrevoSenderEmail: getEnv("BREVO_SENDER_EMAIL", "no-reply@klikku.arjism.com"),

		FrontendURL:   getEnv("FRONTEND_URL", "https://klikku.arjism.com"),
		SuperAdminEmail: getEnv("SUPER_ADMIN_EMAIL", "admin@klikku.arjism.com"),
		SuperAdminPassword: getEnv("SUPER_ADMIN_PASSWORD", "admin123"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func getInt(key string, fallback int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return fallback
}

func getBool(key string, fallback bool) bool {
	if val := os.Getenv(key); val != "" {
		if b, err := strconv.ParseBool(val); err == nil {
			return b
		}
	}
	return fallback
}

func getDuration(key string, fallbackMin time.Duration) time.Duration {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return time.Duration(i) * time.Minute
		}
	}
	return fallbackMin * time.Minute
}

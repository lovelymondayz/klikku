package config

import (
	"fmt"
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

	BrevoAPIKey    string
	BrevoSender    string
	BrevoSenderEmail string

	FrontendURL   string
	SuperAdminEmail string
	SuperAdminPassword string
}

func Load() *Config {
	cfg := &Config{
		AppPort:     getEnv("APP_PORT", "8083"),
		AppEnv:      getEnv("APP_ENV", "production"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
		JWTExpiry:   getDuration("JWT_EXPIRY_MIN", 15),
		RefreshExpiry: getDuration("REFRESH_EXPIRY_HOUR", 168),

		DBHost:     getEnv("DB_HOST", "postgres"),
		DBPort:     getInt("DB_PORT", 5432),
		DBUser:     getEnv("DB_USER", "klikku"),
		DBPassword: getEnv("DB_PASSWORD", "klikku"),
		DBName:     getEnv("DB_NAME", "klikku"),

		BrevoAPIKey:    os.Getenv("BREVO_API_KEY"),
		BrevoSender:    getEnv("BREVO_SENDER", "Klikku"),
		BrevoSenderEmail: getEnv("BREVO_SENDER_EMAIL", "no-reply@klikku.arjism.com"),

		FrontendURL:   getEnv("FRONTEND_URL", "https://klikku.arjism.com"),
		SuperAdminEmail: getEnv("SUPER_ADMIN_EMAIL", "admin@klikku.arjism.com"),
		SuperAdminPassword: os.Getenv("SUPER_ADMIN_PASSWORD"),
	}

	if cfg.JWTSecret == "" {
		panic("JWT_SECRET environment variable is required")
	}

	if cfg.AppEnv == "production" && cfg.SuperAdminPassword == "" {
		panic("SUPER_ADMIN_PASSWORD environment variable is required in production")
	}

	if cfg.SuperAdminPassword == "" {
		cfg.SuperAdminPassword = "admin123"
		fmt.Println("⚠️  WARNING: Using default super admin password. Set SUPER_ADMIN_PASSWORD env var!")
	}

	return cfg
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

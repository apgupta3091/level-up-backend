package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string
	Env  string

	DatabaseURL string

	JWTSecret        string
	JWTAccessExpiry  time.Duration
	JWTRefreshExpiry time.Duration

	StripeSecretKey     string
	StripeWebhookSecret string
	StripePriceID       string

	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
	EmailFrom    string
}

func Load() (*Config, error) {
	// Load .env file if present (ignored in production where env vars are set directly)
	_ = godotenv.Load()

	cfg := &Config{
		Port:        getEnv("PORT", "8080"),
		Env:         getEnv("ENV", "development"),
		DatabaseURL: requireEnv("DATABASE_URL"),

		JWTSecret:        requireEnv("JWT_SECRET"),
		JWTAccessExpiry:  parseDuration("JWT_ACCESS_EXPIRY", 15*time.Minute),
		JWTRefreshExpiry: parseDuration("JWT_REFRESH_EXPIRY", 7*24*time.Hour),

		StripeSecretKey:     requireEnv("STRIPE_SECRET_KEY"),
		StripeWebhookSecret: requireEnv("STRIPE_WEBHOOK_SECRET"),
		StripePriceID:       requireEnv("STRIPE_PRICE_ID"),

		SMTPHost:     getEnv("SMTP_HOST", "localhost"),
		SMTPPort:     parseInt("SMTP_PORT", 1025),
		SMTPUser:     getEnv("SMTP_USER", ""),
		SMTPPassword: getEnv("SMTP_PASSWORD", ""),
		EmailFrom:    getEnv("EMAIL_FROM", "noreply@levelup.dev"),
	}

	return cfg, nil
}

func (c *Config) IsDevelopment() bool {
	return c.Env == "development"
}

func requireEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic(fmt.Sprintf("required environment variable %q is not set", key))
	}
	return v
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func parseDuration(key string, fallback time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return fallback
	}
	return d
}

func parseInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return n
}

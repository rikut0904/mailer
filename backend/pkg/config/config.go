package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port               string
	DatabaseURL        string
	FirebaseProjectID  string
	FirebaseAPIKey     string
	FirebaseAuthDomain string
	AllowedOrigins     []string
	AutoMigrate        bool
}

func Load() (*Config, error) {
	cfg := &Config{
		Port:               getEnv("PORT", "8080"),
		DatabaseURL:        os.Getenv("DATABASE_URL"),
		FirebaseProjectID:  os.Getenv("FIREBASE_PROJECT_ID"),
		FirebaseAPIKey:     os.Getenv("FIREBASE_API_KEY"),
		FirebaseAuthDomain: os.Getenv("FIREBASE_AUTH_DOMAIN"),
		AllowedOrigins:     []string{getEnv("ALLOWED_ORIGIN", "http://localhost:3000")},
		AutoMigrate:        getEnvBool("AUTO_MIGRATE", true),
	}

	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	if v := os.Getenv(key); v != "" {
		switch v {
		case "1", "true", "TRUE", "yes", "YES", "on", "ON":
			return true
		case "0", "false", "FALSE", "no", "NO", "off", "OFF":
			return false
		}
	}
	return fallback
}

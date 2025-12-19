package config

import (
	"os"
)

type Config struct{
	HTTPPort string
	LogLevel string
	DBDSN    string
}

func Load() *Config {
	cfg := &Config{
		HTTPPort: getEnv("HTTP_PORT", "8080"),
		LogLevel: getEnv("LOG_LEVEL", "info"),
		DBDSN:    getEnv("DB_DSN", ""),
	}

	if cfg.DBDSN == "" {
		panic("DB_DSN is required")
	}

	return cfg
}

func getEnv(key, defaultValue string) string {
	v := os.Getenv(key)
	if v != "" {
		return v
	}
	return defaultValue
}

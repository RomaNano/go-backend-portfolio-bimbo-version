package config

import (
	"os"
)

type Config struct{
	HTTPPort string
}

func Load() *Config {
	return &Config{
		HTTPPort: getEnv("HTTP_PORT", "8080")
	}
}

func getEnv(key, defaultValue string) string {
	v := os.Getenv(key)
	if v! = "" {
		return v
	}
	return defaultValue
}

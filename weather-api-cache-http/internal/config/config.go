package config

import (
	"os"
)

type Config struct {
	HTTPPort string
	RedisAddr string
}


func Load() *Config {
	return &Config{
		HTTPPort: getEnv("HTTP_PORT","8080"),
		RedisAddr: getEnv("REDIS_ADDR", "localhost:6379"),
	}
}


func getEnv(key, defaultValue string) string {
	if v:= os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}
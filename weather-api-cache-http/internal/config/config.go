package config

import (
	"os"
	"strconv"
)

type Config struct {
	HTTPPort string
	RedisAddr string
	WeatherBaseURL string
	WeatherTimeoutSec int
	WeatherCacheTTLSeconds int
}


func Load() *Config {
	return &Config{
		HTTPPort:  getEnv("HTTP_PORT", "8080"),
		RedisAddr: getEnv("REDIS_ADDR", "localhost:6379"),

		WeatherBaseURL: getEnv("WEATHER_BASE_URL", "https://api.open-meteo.com/v1/forecast"),
		WeatherTimeoutSec: mustInt(getEnv("WEATHER_TIMEOUT_SEC", "3")),
		WeatherCacheTTLSeconds: mustInt(getEnv("WEATHER_CACHE_TTL_SEC", "300")),
	}
}


func getEnv(key, defaultValue string) string {
	if v:= os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

func mustInt(v string) int {
	n, err := strconv.Atoi(v)
	if err != nil {
		panic("invalid int config value: " + v)
	}
	return n
}
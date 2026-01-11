package main

import (
	"log"
	"net/http"
	"time"

	//"weather-api-cache-http/internal/cache/redis"
	"weather-api-cache-http/internal/config"
	"weather-api-cache-http/internal/handler"
	"weather-api-cache-http/internal/logger"
	"weather-api-cache-http/internal/middleware"
	redisCache "weather-api-cache-http/internal/cache/redis"
	weatherClient "weather-api-cache-http/internal/client/weather"
	weatherService "weather-api-cache-http/internal/service/weather"
)

func main() {
	cfg := config.Load()

	logg := logger.New()
	cache := redisCache.New(cfg.RedisAddr)

	apiClient :=weatherClient.New(
		cfg.WeatherBaseURL,
		time.Duration(cfg.WeatherTimeoutSec)*time.Second,
	)



	svc := weatherService.New(
		cache,
		apiClient,
		time.Duration(cfg.WeatherCacheTTLSeconds)*time.Second,
		logg,
	)
	

	mux := http.NewServeMux()
	mux.Handle("/health", handler.Health())
	mux.Handle("/weather", handler.Weather(svc))

	// оборачиваем весь mux, а не конкретный handler
	var h http.Handler = mux
	h = middleware.RequestID(h)
	h = middleware.Logging(logg)(h)

	server := &http.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: h,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout: 60 * time.Second,
	}

	log.Printf("starting http server on :%s", cfg.HTTPPort)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

}

package main

import (
	"log"
	"net/http"

	"weather-api-cache-http/internal/cache/redis"
	"weather-api-cache-http/internal/config"
	"weather-api-cache-http/internal/handler"
	"weather-api-cache-http/internal/logger"
	"weather-api-cache-http/internal/middleware"
)

func main() {
	cfg := config.Load()
	logg := logger.New()
	redisCache := redis.New(cfg.RedisAddr)
	_ = redisCache // временно, уберём на следующем шаге

	mux := http.NewServeMux()
	mux.Handle("/health", handler.Health())

	// оборачиваем весь mux, а не конкретный handler
	var h http.Handler = mux
	h = middleware.RequestID(h)
	h = middleware.Logging(logg)(h)

	server := &http.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: h,
	}

	log.Printf("starting http server on :%s", cfg.HTTPPort)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

}

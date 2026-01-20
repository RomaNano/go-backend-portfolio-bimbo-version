package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	redisCache "weather-api-cache-http/internal/cache/redis"
	weatherClient "weather-api-cache-http/internal/client/weather"
	"weather-api-cache-http/internal/config"
	"weather-api-cache-http/internal/handler"
	"weather-api-cache-http/internal/logger"
	"weather-api-cache-http/internal/metrics"
	"weather-api-cache-http/internal/middleware"
	weatherService "weather-api-cache-http/internal/service/weather"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	cfg := config.Load()

	logg := logger.New()
	cache := redisCache.New(cfg.RedisAddr)

	apiClient := weatherClient.New(
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

	metrics.Register()
	mux.Handle("/metrics", promhttp.Handler())

	// оборачиваем весь mux, а не конкретный handler
	var h http.Handler = mux
	h = middleware.RequestID(h)
	h = middleware.Logging(logg)(h)

	server := &http.Server{
		Addr:              ":" + cfg.HTTPPort,
		Handler:           h,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// context, который отменится при SIGINT / SIGTERM
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	go func() {
		log.Printf("starting http server on: %s", cfg.HTTPPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()

	log.Println("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("server shutdown failed: %v", err)
	}

	log.Println("server gracefully stopped")

}

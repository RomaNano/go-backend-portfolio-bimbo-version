# Weather API with Cache, Concurrency and Observability

A production-style HTTP service written in Go that provides current weather data
with Redis caching, concurrent external API calls and full observability
(logging, metrics, graceful shutdown).


This project is part of a backend portfolio and demonstrates real-world Go
backend patterns.

---

## Features

- Pure `net/http` server (no frameworks)
- External Weather API integration (Open-Meteo)
- Redis cache (cache-aside pattern)
- Concurrent external requests using `errgroup`
- Background cache writes with context handling
- Structured logging (`log/slog`)
- Prometheus metrics (`/metrics`)
- Graceful shutdown (`SIGINT`, `SIGTERM`)
- Configuration via environment variables
- Business logic is covered by unit tests; basic API handler test is included using httptest.

---

## Architecture Overview
HTTP
└── handler
└── service
├── cache (Redis)
└── external API (Open-Meteo)


Key architectural decisions:
- **Service layer** contains business logic
- **Cache is best-effort** and does not block responses
- **Concurrency** is used for external API calls
- **Observability first**: logs + metrics by default

---

##  Configuration

The service is configured via environment variables:

| Variable | Description | Default |
|--------|------------|---------|
| `HTTP_PORT` | HTTP server port | `8080` |
| `REDIS_ADDR` | Redis address | `localhost:6379` |
| `WEATHER_CACHE_TTL_SECONDS` | Cache TTL | `300` |

Example:

```bash
export HTTP_PORT=8080
export REDIS_ADDR=localhost:6379
export WEATHER_CACHE_TTL_SECONDS=300

# 1. Start Redis
docker run -p 6379:6379 redis

# 2. Run the service
go run ./cmd/api

# Result
starting http server on :8080


# API Endpoints
# Get Weather
curl "http://localhost:8080/weather?lat=52.52&lon=13.41"

# Example response:
{
  "temperature": -3.1,
  "windspeed": 13.2,
  "hourly_temperatures": [ ... ]
}

# Metrics (Prometheus)
curl http://localhost:8080/metrics

#Health Check
curl http://localhost:8080/health
```

## Observability
Logging
- Structured logs (slog)
- Request latency
- Cache hit / miss
- Error diagnostics

Metrics
- weather_requests_total
- weather_cache_hits_total
- weather_cache_misses_total
- weather_request_duration_seconds

## Key Concepts Demonstrated
This project intentionally focuses on backend fundamentals:
- Cache-aside pattern
- Context cancellation handling
- Background operations with timeout
- Concurrent API calls (errgroup)
- Graceful shutdown
- Production-grade observability

## Tech Stack
- Go
- net/http
- Redis
- Prometheus
- Open-Meteo API



## Notes
This project is intentionally minimal in scope but deep in implementation,
focusing on correctness, clarity and production-ready patterns rather than
feature completeness.
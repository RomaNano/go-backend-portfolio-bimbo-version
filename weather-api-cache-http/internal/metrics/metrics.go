package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	WeatherRequests = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "weather_requests_total",
			Help: "Total number of weather requests",
		},
	)

	CacheHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "weather_cache_hits_total",
			Help: "Total number of weather cache hits",
		},
	)

	CacheMisses = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "weather_cache_misses_total",
			Help: "Total number of weather cache misses",
		},
	)

	RequestDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name: "weather_request_duration_seconds",
			Help: "Weather request latency",
			Buckets: prometheus.DefBuckets,
		},
	)
)

func Register() {
	prometheus.MustRegister(
		WeatherRequests,
		CacheHits,
		CacheMisses,
		RequestDuration,
	)
}
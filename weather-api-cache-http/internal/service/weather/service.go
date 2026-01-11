package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"log/slog"

	"weather-api-cache-http/internal/cache"
	"weather-api-cache-http/internal/client/weather"
	"weather-api-cache-http/internal/model"	
	"weather-api-cache-http/internal/metrics"

	"golang.org/x/sync/errgroup"
)


type Service struct {
	cache cache.Cache
	api weather.Client
	ttl time.Duration
	log *slog.Logger
}

func New(cache cache.Cache, api weather.Client, ttl time.Duration, log *slog.Logger) *Service {
	return &Service{cache: cache, api: api, ttl: ttl, log: log}
}

func cacheKey(lat, lon float64) string {
	return fmt.Sprintf(
		"weather:current:%d:%d",
		int(lat*10000),
		int(lon*10000),
	)
}

func (s *Service) GetCurrent(
	ctx context.Context,
	lat, lon float64,
) (*model.Weather, error) {

	key := cacheKey(lat, lon)

	// 1) try cache
	cached, err := s.cache.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	s.log.Info(
		"cache lookup",
		"key", key,
		"cached_len", len(cached),
	)

	if cached != "" {
		metrics.CacheHits.Inc()
		s.log.Info("weather cache hit", "key", key)

		var w model.Weather
		if err := json.Unmarshal([]byte(cached), &w); err == nil {
			return &w, nil
		}

		// если кеш битый — считаем как miss и идём дальше
		s.log.Warn("failed to unmarshal cache, fallback to api", "key", key)
	}

	// cache miss (или битый кеш)
	metrics.CacheMisses.Inc()
	s.log.Info("weather cache miss", "key", key)

	// 2) concurrent external calls
	var (
		current *model.Weather
		hourly  []float64
	)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		w, err := s.api.GetCurrent(ctx, lat, lon)
		if err != nil {
			return err
		}
		current = w
		return nil
	})

	g.Go(func() error {
		h, err := s.api.GetHourly(ctx, lat, lon)
		if err != nil {
			return err
		}
		hourly = h
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	current.HourlyTemps = hourly

	// 3) save to cache
	b, err := json.Marshal(current)
	if err == nil {
		cacheCtx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
		defer cancel()

		if err := s.cache.Set(cacheCtx, key, string(b), s.ttl); err != nil {
			s.log.Warn("failed to set cache", "err", err, "key", key)
		}
	}

	return current, nil
}




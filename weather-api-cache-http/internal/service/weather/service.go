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
	return fmt.Sprintf("weather:current:%0.4f:%0.4f", lat, lon)
}

func (s *Service) GetCurrent(ctx context.Context, lat, lon float64) (*model.Weather, error) {

	key := cacheKey(lat, lon)

	// 1) cache
	cached, err := s.cache.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	if cached != "" {
		s.log.Info("weather cache hit", "key", key)
		var w model.Weather
		if err := json.Unmarshal([]byte(cached), &w); err == nil {
			return &w, nil
		}
	} else {
		s.log.Info("weather cache miss", "key", key)
	}

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
		_ = s.cache.Set(ctx, key, string(b), s.ttl)
	}

	return current, nil
}



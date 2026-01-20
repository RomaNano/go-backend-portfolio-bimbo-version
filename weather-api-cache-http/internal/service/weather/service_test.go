package weather

import (
	"context"
	"io"
	"log/slog"
	"testing"
	"time"

	"weather-api-cache-http/internal/model"
)

// Fake implementations

type fakeCache struct {
	data map[string]string
}

func newFakeCache() *fakeCache {
	return &fakeCache{data: make(map[string]string)}
}

func (f *fakeCache) Get(ctx context.Context, key string) (string, error) {
	return f.data[key], nil
}

func (f *fakeCache) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	f.data[key] = value
	return nil
}

// Fake weather API

type fakeWeatherAPI struct {
	currentCalls int
	hourlyCalls  int
}

// Mocking
func (f *fakeWeatherAPI) GetCurrent(
	ctx context.Context,
	lat, lon float64,
) (*model.Weather, error) {
	f.currentCalls++
	return &model.Weather{
		Temperature: 10,
		WindSpeed:   5,
	}, nil
}

func (f *fakeWeatherAPI) GetHourly(
	ctx context.Context,
	lat, lon float64,
) ([]float64, error) {
	f.hourlyCalls++
	return []float64{1, 2, 3}, nil
}

// Тест: cache miss → API вызывается

func TestGetCurrent_CacheMiss(t *testing.T) {

	cache := newFakeCache()
	api := &fakeWeatherAPI{}

	svc := New(
		cache,
		api,
		5*time.Minute,
		slog.New(slog.NewTextHandler(io.Discard, nil)),
	)

	ctx := context.Background()

	weather, err := svc.GetCurrent(ctx, 52.52, 13.41)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if weather == nil {
		t.Fatal("expected weather, got nil")
	}

	if api.currentCalls != 1 || api.hourlyCalls != 1 {
		t.Fatalf("expected api to be called, got current=%d hourly=%d",
			api.currentCalls, api.hourlyCalls)
	}
}

// Тест: cache hit → API НЕ вызывается

func TestGetCurrent_CacheHit(t *testing.T) {
	cache := newFakeCache()
	api := &fakeWeatherAPI{}

	// кладём данные в кеш заранее
	cached := `{"temperature":10,"windspeed":5,"hourly_temperatures":[1,2,3]}`
	cache.data[cacheKey(52.52, 13.41)] = cached

	svc := New(
		cache,
		api,
		5*time.Minute,
		slog.New(slog.NewTextHandler(io.Discard, nil)),
	)

	ctx := context.Background()

	weather, err := svc.GetCurrent(ctx, 52.52, 13.41)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if weather == nil {
		t.Fatal("expected weather, got nil")
	}

	if api.currentCalls != 0 || api.hourlyCalls != 0 {
		t.Fatalf("expected api NOT to be called, got current=%d hourly=%d",
			api.currentCalls, api.hourlyCalls)
	}
}

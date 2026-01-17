package handler_test

import (
	"context"
	"time"
	"io"
	"log/slog"
	"net/http"
	"testing"
	"net/http/httptest"
	"encoding/json"

	"weather-api-cache-http/internal/model"	
	"weather-api-cache-http/internal/service/weather"
	"weather-api-cache-http/internal/handler"
)

// Готовим fake-зависимости

type fakeCache struct {}

func (f *fakeCache) Get(ctx context.Context, key string) (string, error){
	return "", nil
}

func (f *fakeCache) Set(ctx context.Context, key, value string, ttl time.Duration) error{
	return nil
}

// Готовим fake Weather API

type fakeWeatherAPI struct{}

func (f *fakeWeatherAPI) GetCurrent(
	ctx context.Context,
	lat, lon float64,
) (*model.Weather, error) {
	return &model.Weather{
		Temperature: 10,
		WindSpeed:   5,
	}, nil
}

func (f *fakeWeatherAPI) GetHourly(
	ctx context.Context,
	lat, lon float64,
) ([]float64, error) {
	return []float64{1,2,3},nil
}

// Собираем сервис и handler

func setupTestServer() http.Handler {
	cache := &fakeCache{}
	api := &fakeWeatherAPI{}

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	svc := weather.New(
		cache,
		api,
		time.Minute,
		logger,
	)

	// ????
	h := handler.Weather(svc)

	mux := http.NewServeMux()
	mux.Handle("/weather", h)

	return mux
}


// Сам API-тест

func TestWeatherAPI_OK(t *testing.T) {
	server := setupTestServer()

	req := httptest.NewRequest(
		http.MethodGet,
		"/weather?lat=52.52&lon=13.41",
		nil,
	)

	rec := httptest.NewRecorder()

	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var resp map[string]interface{}

	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid json response: %v", err)
	}

	if resp["temperature"] == nil {
		t.Fatalf("expected temperature field")
	}

}
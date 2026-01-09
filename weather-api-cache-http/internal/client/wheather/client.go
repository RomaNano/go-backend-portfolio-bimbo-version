package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"weather-api-cache-http/internal/model"
)

type Client interface {
	GetCurrent(ctx context.Context, lat, lon float64)(*model.Weather, err)
}

type HTTPClient struct {
	baseURL string
	client *http.Client
}

func New(baseURL string, timeout time.Duration) *HTTPClient {
	return &HTTPClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

type apiResponse struct{
	CurrentWeather struct{
		Temperature float64 `json:"temperature"`
		WindSpeed   float64 `json:"windspeed"`
	} `json:"current_weather"`
}

func (c *HTTPClient) GetCurrent(
	ctx context.Context,
	lat, lon float64,
) (*model.Weather, error) {

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf(
			"%s?latitude=%f&longitude=%f&current_weather=true",
			c.baseURL,
			lat,
			lon,
		),
		nil,
	)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("weather api returned %s", resp.Status)
	}

	var apiResp apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	return &model.Weather{
		Temperature: apiResp.CurrentWeather.Temperature,
		WindSpeed:   apiResp.CurrentWeather.WindSpeed,
	}, nil
}


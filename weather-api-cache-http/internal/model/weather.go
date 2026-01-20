package model

// DTO
type Weather struct {
	Temperature float64   `json:"temperature"`
	WindSpeed   float64   `json:"windspeed"`
	HourlyTemps []float64 `json:"hourly_temperatures"`
}

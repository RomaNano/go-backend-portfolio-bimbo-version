package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"weather-api-cache-http/internal/service/weather"
	"weather-api-cache-http/internal/metrics"
	"github.com/prometheus/client_golang/prometheus"

)

func Weather(svc *weather.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
		metrics.WeatherRequests.Inc()
		timer := prometheus.NewTimer(metrics.RequestDuration)
		defer timer.ObserveDuration()
		
		q := r.URL.Query()

		latStr := q.Get("lat")
		lonStr := q.Get("lon")

		lat, err := strconv.ParseFloat(latStr, 64)
		if err != nil {
			http.Error(w, "invalid lat", http.StatusBadRequest)
			return
		}

		lon, err := strconv.ParseFloat(lonStr, 64)
		if err != nil {
			http.Error(w, "invalid lon", http.StatusBadRequest)
			return
		}

		res, err := svc.GetCurrent(r.Context(),lat,lon)
		if err != nil {
			http.Error(w, "failed to get weather", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(res)
	})
}
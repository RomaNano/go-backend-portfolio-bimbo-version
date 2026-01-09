package middleware

import (
	"net/http"
	"time"

	"log/slog"
)

func Logging(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			next.ServeHTTP(w, r)

			duration := time.Since(start)

			reqID := r.Context().Value(requestIDKey)
			
		log.Info("http request",
			"method", r.Method,
			"path", r.URL.Path,
			"request_id", reqID,
			"duration", duration.String(),
		)
		})
	}

}
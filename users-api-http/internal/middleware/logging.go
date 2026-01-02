package middleware

import (
	"net/http"
	"time"
	"log"
)

// перехватчик, тк мы не нет способа узнать 
// статус ответа после того, как handler отработал.
// Текущая обёртка запомнит status code и запроксирует вызовы дальше
type responseWriter struct{
	http.ResponseWriter 
	status int
}

func (rw responseWriter) WriteHeader(code int){
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {

	// если handler не вызвал WriteHeader — статус будет 200
	if rw.status == 0 {
		rw.status= http.StatusOK
	}

	return rw.ResponseWriter.Write(b)

}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &responseWriter{
			ResponseWriter: w,
		}

		next.ServeHTTP(rw, r)
		duration := time.Since(start)
		reqID := GetRequestID(r.Context())

		log.Printf(
			"method=%s path=%s status=%d duration=%s request_id=%s",
			r.Method,
			r.URL.Path,
			rw.status,
			duration,
			reqID,
		)

	})
}
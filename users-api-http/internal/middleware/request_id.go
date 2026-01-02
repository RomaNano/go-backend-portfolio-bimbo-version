package middleware


import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type contextKey string 

const requestIDKey contextKey = "request_id"

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		reqID := r.Header.Get("X-Request-Id") 
		if reqID == "" {
			reqID = uuid.NewString()
		}

		w.Header().Set("X-Request-Id", reqID)
		ctx := context.WithValue(r.Context(), requestIDKey, reqID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetRequestID(ctx context.Context) string{
	if v := ctx.Value(requestIDKey); v!=nil {
		if id, ok := v.(string); ok {
			return id
		}
	}
	return ""
}
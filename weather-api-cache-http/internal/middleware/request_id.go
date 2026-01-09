package middleware

import (
	"net/http"
	"context"

	"github.com/google/uuid"
)


type ctxKey string

const requestIDKey ctxKey = "request_id"
const requestIDHeader = "X-Request-ID"

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		reqID := r.Header.Get(requestIDHeader)
		if reqID == ""{
			reqID = uuid.NewString()
		}	
	
	// суём в context
	ctx := context.WithValue(r.Context(), requestIDKey, reqID)

	// добавляем в response
	w.Header().Set(requestIDHeader, reqID)

	// передаём дальше
	next.ServeHTTP(w, r.WithContext(ctx))

	})
}


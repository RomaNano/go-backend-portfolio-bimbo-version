package handler

import (
	"net/http"
)

func Health() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
}
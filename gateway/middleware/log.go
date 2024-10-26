package middleware

import (
	"log"
	"net/http"
	"time"
)

func Log() Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			h.ServeHTTP(w, r)
			log.Printf("%s %s %v\n", r.Method, r.URL.Path, time.Since(start))
		})
	}
}

package middleware

import (
	"log"
	"net/http"
	"time"
)

// function yang menerima handler, return handler juga
func Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		log.Printf(
			"[REQUEST] %s %s dari %s",
			r.Method,
			r.RequestURI,
			r.RemoteAddr,
		)

		next(w, r)

		duration := time.Since(start)
		log.Printf(
			"[DONE] %s %s selesai dalam %v",
			r.Method,
			r.RequestURI,
			duration,
		)
	}
}

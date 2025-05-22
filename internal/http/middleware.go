package http

import (
	"k8s-mipt/internal/metrics"
	"net/http"
	"time"
)

func durationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		duration := time.Since(start).Seconds()
		metrics.RequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration)
	})
}

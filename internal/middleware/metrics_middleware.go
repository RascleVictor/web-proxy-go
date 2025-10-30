package middleware

import (
	"net/http"
	"time"

	"web-proxy-go/internal/metrics"
)

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(rec, r)
		duration := time.Since(start).Seconds()

		backend := r.Host
		metrics.ProxyLatency.WithLabelValues(backend).Observe(duration)
		if rec.statusCode >= 400 {
			metrics.ProxyErrors.WithLabelValues(http.StatusText(rec.statusCode)).Inc()
		}
	})
}

type metricsResponseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (r *metricsResponseRecorder) WriteHeader(code int) {
	r.statusCode = code
	r.ResponseWriter.WriteHeader(code)
}

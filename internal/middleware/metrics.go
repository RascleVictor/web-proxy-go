package middleware

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	httpRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "proxy_http_requests_total",
			Help: "Nombre total de requêtes HTTP reçues",
		},
		[]string{"method", "path"},
	)

	httpDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "proxy_http_duration_seconds",
			Help:    "Durée des requêtes HTTP",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)
)

func init() {
	prometheus.MustRegister(httpRequests, httpDuration)
}

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		duration := time.Since(start).Seconds()
		httpRequests.WithLabelValues(r.Method, r.URL.Path).Inc()
		httpDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration)
	})
}

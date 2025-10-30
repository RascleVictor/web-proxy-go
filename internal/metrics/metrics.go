package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	ProxyLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "proxy_latency_seconds",
			Help:    "Latence des requêtes proxy par backend",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"backend"},
	)

	ProxyErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "proxy_errors_total",
			Help: "Nombre d'erreurs HTTP par code",
		},
		[]string{"code"},
	)
)

func Init() {
	prometheus.MustRegister(ProxyLatency, ProxyErrors)
}

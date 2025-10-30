package proxy

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"
	"time"

	"go.uber.org/zap"
	"golang.org/x/net/http2"

	"web-proxy-go/internal/logger"
	"web-proxy-go/internal/metrics"

	"go.opentelemetry.io/otel"
)

// --- Load Balancer ---
type LoadBalancer struct {
	backends []*url.URL
	index    uint64
}

func NewLoadBalancer(backendURLs []string) *LoadBalancer {
	backends := make([]*url.URL, len(backendURLs))
	for i, b := range backendURLs {
		u, err := url.Parse(b)
		if err != nil {
			logger.Log.Fatal("Invalid backend URL", zap.String("backend", b), zap.Error(err))
		}
		backends[i] = u
	}
	return &LoadBalancer{backends: backends}
}

func (lb *LoadBalancer) NextBackend() *url.URL {
	i := atomic.AddUint64(&lb.index, 1)
	return lb.backends[i%uint64(len(lb.backends))]
}

// --- Reverse Proxy ---
func NewProxy(lb *LoadBalancer) *httputil.ReverseProxy {
	transport := &http.Transport{
		MaxIdleConns:          200,
		MaxIdleConnsPerHost:   50,
		IdleConnTimeout:       90 * time.Second,
		DisableCompression:    false,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
	}

	_ = http2.ConfigureTransport(transport)

	return &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			target := lb.NextBackend()
			req.URL.Scheme = target.Scheme
			req.URL.Host = target.Host
			req.Host = target.Host

			// Timeout global
			ctx, cancel := context.WithTimeout(req.Context(), 15*time.Second)
			defer cancel()

			// Tracing OpenTelemetry
			tracer := otel.Tracer("web-proxy")
			ctx, span := tracer.Start(ctx, "proxy-request")
			defer span.End()

			// On ajoute seulement la startTime dans le contexte (pour mesurer la durée)
			req = req.WithContext(contextWithStartTime(ctx))
		},
		Transport: transport,
		ModifyResponse: func(res *http.Response) error {
			start := getStartTime(res.Request.Context())
			duration := time.Since(start)

			// Récupération du RequestID si le middleware l'a injecté
			reqID := ""
			if v := res.Request.Context().Value("requestID"); v != nil {
				if id, ok := v.(string); ok {
					reqID = id
				}
			}

			backend := res.Request.URL.Host

			// Logs enrichis
			logger.Log.Info("Requête proxy terminée",
				zap.String("requestID", reqID),
				zap.String("backend", backend),
				zap.String("url", res.Request.URL.String()),
				zap.Int("status", res.StatusCode),
				zap.Duration("duration", duration),
				zap.Int("responseSize", int(res.ContentLength)),
			)

			// Metrics Prometheus
			metrics.ProxyLatency.WithLabelValues(backend).Observe(duration.Seconds())
			if res.StatusCode >= 400 {
				metrics.ProxyErrors.WithLabelValues(fmt.Sprint(res.StatusCode)).Inc()
			}

			return nil
		},
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			reqID := ""
			if v := r.Context().Value("requestID"); v != nil {
				if id, ok := v.(string); ok {
					reqID = id
				}
			}
			logger.Log.Error("Erreur proxy",
				zap.String("requestID", reqID),
				zap.Error(err),
			)
			http.Error(w, "Erreur proxy : "+err.Error(), http.StatusBadGateway)
		},
	}
}

// --- Context pour calculer durée ---
type ctxKey string

const startTimeKey ctxKey = "startTime"

func contextWithStartTime(ctx context.Context) context.Context {
	return context.WithValue(ctx, startTimeKey, time.Now())
}

func getStartTime(ctx context.Context) time.Time {
	if v := ctx.Value(startTimeKey); v != nil {
		if t, ok := v.(time.Time); ok {
			return t
		}
	}
	return time.Now()
}

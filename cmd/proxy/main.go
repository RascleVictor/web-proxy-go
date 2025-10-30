package main

import (
	"fmt"
	"net/http"

	"web-proxy-go/internal/config"
	"web-proxy-go/internal/logger"
	"web-proxy-go/internal/middleware"
	"web-proxy-go/internal/proxy"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func main() {
	cfg := config.LoadConfig()
	logger.InitLogger(cfg.Logging.Level)
	defer logger.Log.Sync()

	// LoadBalancer
	lb := proxy.NewLoadBalancer(cfg.Server.Backends)

	// Création du proxy (plus d'erreur, une seule valeur)
	proxyHandler := proxy.NewProxy(lb)

	// Chaîne de middlewares
	handler := middleware.RecoveryMiddleware(
		middleware.LoggingMiddleware(
			middleware.MetricsMiddleware(
				middleware.CORSMiddleware(proxyHandler),
			),
		),
	)

	http.Handle("/", handler)
	http.Handle("/metrics", promhttp.Handler())

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	logger.Log.Info("Proxy démarré", zap.String("addr", addr))

	if err := http.ListenAndServe(addr, nil); err != nil {
		logger.Log.Fatal("Erreur serveur", zap.Error(err))
	}
}

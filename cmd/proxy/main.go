package main

import (
	"fmt"
	"net/http"
	"web-proxy-go/internal/config"
	"web-proxy-go/internal/logger"
	"web-proxy-go/internal/metrics"
	"web-proxy-go/internal/middleware"
	"web-proxy-go/internal/proxy"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func main() {
	cfg := config.LoadConfig()
	logger.InitLogger(cfg.Logging.Level)
	defer logger.Log.Sync()

	metrics.Init()

	lb := proxy.NewLoadBalancer(cfg.Server.Backends)

	proxyHandler := proxy.NewProxy(lb)

	ipFilter := middleware.NewIPFilter(
		[]string{"192.168.1.10"}, // whitelist
		[]string{"10.0.0.5"},     // blacklist
	)

	// Compose les middlewares : RequestID en premier (au plus tôt), puis recovery, logging, metrics, CORS
	handler := middleware.RequestIDMiddleware(
		ipFilter.Middleware(

			middleware.RecoveryMiddleware(
				middleware.LoggingMiddleware(
					middleware.MetricsMiddleware(
						middleware.CORSMiddleware(proxyHandler),
					),
				),
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

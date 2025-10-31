package main

import (
	"fmt"
	"net/http"
	"time"
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

	ipFilter := middleware.NewDynamicIPFilter(
		[]string{"127.0.0.1"},
		100,
		10,
		1*time.Minute,
		5*time.Minute,
		30*time.Second,
	)

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

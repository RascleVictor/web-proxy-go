package main

import (
	"fmt"
	"net/http"

	"proxy-web-go/internal/config"
	"proxy-web-go/internal/logger"
	"proxy-web-go/internal/middleware"
	"proxy-web-go/internal/proxy"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func main() {
	cfg := config.LoadConfig()
	logger.InitLogger(cfg.Logging.Level)
	defer logger.Log.Sync()

	proxyHandler, err := proxy.NewProxy(cfg.Server.Target)
	if err != nil {
		logger.Log.Fatal("Erreur création proxy", zap.Error(err))
	}

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

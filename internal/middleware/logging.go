package middleware

import (
	"go.uber.org/zap"
	"net/http"
	"time"
	"web-proxy-go/internal/logger"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)

		logger.Log.Info("request okay",
			zap.String("method", r.Method),
			zap.String("url", r.URL.String()),
			zap.String("remote_addr", r.RemoteAddr),
			zap.Duration("duration", time.Since(start)))
	})
}

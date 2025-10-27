package middleware

import (
	"go.uber.org/zap"
	"net/http"
	"proxy-web-go/internal/logger"
)

func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.Log.Error("pPanic intercepter", zap.Any("err", err))
				http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

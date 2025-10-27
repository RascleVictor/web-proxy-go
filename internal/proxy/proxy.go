package proxy

import (
	"context"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"go.uber.org/zap"
	"proxy-web-go/internal/logger"
)

func NewProxy(target string) (*httputil.ReverseProxy, error) {
	urlTarget, err := url.Parse(target)
	if err != nil {
		return nil, err
	}

	transport := &http.Transport{
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   20,
		IdleConnTimeout:       90 * time.Second,
		DisableCompression:    false,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
	}

	proxy := httputil.NewSingleHostReverseProxy(urlTarget)
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.URL.Scheme = urlTarget.Scheme
		req.URL.Host = urlTarget.Host
		req.Host = urlTarget.Host
		ctx := req.Context()
		ctx = contextWithStartTime(ctx)
		req = req.WithContext(ctx)
	}

	proxy.Transport = transport

	proxy.ModifyResponse = func(res *http.Response) error {
		start := getStartTime(res.Request.Context())
		duration := time.Since(start)
		logger.Log.Info("requête proxy terminée",
			zap.String("url", res.Request.URL.String()),
			zap.Int("status", res.StatusCode),
			zap.Duration("duration", duration),
		)
		return nil
	}

	proxy.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, err error) {
		logger.Log.Error("erreur proxy", zap.Error(err))
		http.Error(writer, "Erreur proxy : "+err.Error(), http.StatusBadGateway)
	}

	return proxy, nil
}

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

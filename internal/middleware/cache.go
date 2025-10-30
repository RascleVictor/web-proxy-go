package middleware

import (
	"bytes"
	"net/http"
	"sync"
	"time"

	lru "github.com/hashicorp/golang-lru"
)

type CachedResponse struct {
	Body       []byte
	Header     http.Header
	StatusCode int
	Timestamp  time.Time
}

var (
	cache     *lru.Cache
	cacheOnce sync.Once
	cacheTTL  = 10 * time.Second // TTL par défaut
)

// Initialiser le cache
func InitCache(size int, ttl time.Duration) {
	cacheOnce.Do(func() {
		c, _ := lru.New(size)
		cache = c
		cacheTTL = ttl
	})
}

// Middleware cache
func CacheMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if cache == nil {
			next.ServeHTTP(w, r)
			return
		}

		key := r.URL.String()
		if val, ok := cache.Get(key); ok {
			cached := val.(CachedResponse)
			if time.Since(cached.Timestamp) < cacheTTL {
				for k, vv := range cached.Header {
					for _, v := range vv {
						w.Header().Add(k, v)
					}
				}
				w.WriteHeader(cached.StatusCode)
				w.Write(cached.Body)
				return
			}
			cache.Remove(key)
		}

		// Capture la réponse
		rec := &responseRecorder{
			ResponseWriter: w,
			body:           &bytes.Buffer{},
		}
		next.ServeHTTP(rec, r)

		cached := CachedResponse{
			Body:       rec.body.Bytes(),
			Header:     rec.Header(),
			StatusCode: rec.statusCode,
			Timestamp:  time.Now(),
		}
		cache.Add(key, cached)
	})
}

type responseRecorder struct {
	http.ResponseWriter
	body       *bytes.Buffer
	statusCode int
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

func (r *responseRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

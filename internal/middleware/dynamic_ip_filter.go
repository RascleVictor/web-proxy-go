package middleware

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

type DynamicIPFilter struct {
	mu sync.Mutex

	whitelist map[string]bool
	banlist   map[string]time.Time
	activity  map[string]*ipStats

	requestLimit    int
	errorLimit      int
	windowDuration  time.Duration
	banDuration     time.Duration
	cleanupInterval time.Duration
}

type ipStats struct {
	requests []time.Time
	errors   int
}

// --- Constructeur ---
func NewDynamicIPFilter(whitelist []string, reqLimit, errLimit int, window, banDuration, cleanup time.Duration) *DynamicIPFilter {
	d := &DynamicIPFilter{
		whitelist:       make(map[string]bool),
		banlist:         make(map[string]time.Time),
		activity:        make(map[string]*ipStats),
		requestLimit:    reqLimit,
		errorLimit:      errLimit,
		windowDuration:  window,
		banDuration:     banDuration,
		cleanupInterval: cleanup,
	}

	for _, ip := range whitelist {
		d.whitelist[ip] = true
	}

	go d.cleanupLoop()
	return d
}

// --- Middleware principal ---
func (d *DynamicIPFilter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := getClientIP(r)

		if d.whitelist[ip] {
			next.ServeHTTP(w, r)
			return
		}

		if d.isBanned(ip) {
			http.Error(w, "Forbidden - IP temporarily banned", http.StatusForbidden)
			return
		}

		d.recordRequest(ip)

		rec := &ipResponseRecorder{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(rec, r)

		if rec.statusCode >= 400 {
			d.recordError(ip)
		}
	})
}

// --- Logique de détection ---
func (d *DynamicIPFilter) recordRequest(ip string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	stats, exists := d.activity[ip]
	if !exists {
		stats = &ipStats{}
		d.activity[ip] = stats
	}
	now := time.Now()

	cutoff := now.Add(-d.windowDuration)
	filtered := []time.Time{}
	for _, t := range stats.requests {
		if t.After(cutoff) {
			filtered = append(filtered, t)
		}
	}
	stats.requests = append(filtered, now)

	if len(stats.requests) >= d.requestLimit {
		d.banlist[ip] = now.Add(d.banDuration)
		delete(d.activity, ip)
	}
}

func (d *DynamicIPFilter) recordError(ip string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	stats, exists := d.activity[ip]
	if !exists {
		stats = &ipStats{}
		d.activity[ip] = stats
	}
	stats.errors++
	if stats.errors >= d.errorLimit {
		d.banlist[ip] = time.Now().Add(d.banDuration)
		delete(d.activity, ip)
	}
}

func (d *DynamicIPFilter) isBanned(ip string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	until, banned := d.banlist[ip]
	if !banned {
		return false
	}
	if time.Now().After(until) {
		delete(d.banlist, ip)
		return false
	}
	return true
}

func (d *DynamicIPFilter) cleanupLoop() {
	for {
		time.Sleep(d.cleanupInterval)
		d.mu.Lock()
		for ip, until := range d.banlist {
			if time.Now().After(until) {
				delete(d.banlist, ip)
			}
		}
		d.mu.Unlock()
	}
}

// --- Utils ---
func getClientIP(r *http.Request) string {
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}

// --- Recorder spécifique à ce middleware ---
type ipResponseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (r *ipResponseRecorder) WriteHeader(status int) {
	r.statusCode = status
	r.ResponseWriter.WriteHeader(status)
}

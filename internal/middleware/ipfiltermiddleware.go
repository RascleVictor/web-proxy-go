package middleware

import (
	"net"
	"net/http"
	"strings"
)

type IPFilter struct {
	Whitelist map[string]bool
	Blacklist map[string]bool
}

func NewIPFilter(whitelist, blacklist []string) *IPFilter {
	wl := make(map[string]bool)
	for _, ip := range whitelist {
		wl[ip] = true
	}
	bl := make(map[string]bool)
	for _, ip := range blacklist {
		bl[ip] = true
	}
	return &IPFilter{
		Whitelist: wl,
		Blacklist: bl,
	}
}

func (f *IPFilter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := getIP(r)

		// Vérifie blacklist en priorité
		if f.Blacklist[ip] {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// Si whitelist définie, vérifier qu'elle inclut l'IP
		if len(f.Whitelist) > 0 && !f.Whitelist[ip] {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Récupère l’IP réelle du client (en prenant X-Forwarded-For en compte)
func getIP(r *http.Request) string {
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}

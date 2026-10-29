package httpapi

import (
	"net"
	"net/http"
	"strings"
)

const contentSecurityPolicy = "default-src 'self'; img-src 'self' https: data:; media-src 'self'; " +
	"frame-ancestors 'none'; base-uri 'none'; form-action 'self'"

// guard protects the local API from other web pages. Browsers may only
// send writes from the dashboard's own origin, and the Host header must be
// an IP address or localhost, which defeats DNS rebinding.
func guard(next http.Handler) http.Handler {
	protected := http.NewCrossOriginProtection().Handler(next)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !localHost(r.Host) {
			http.Error(w, "grabb only answers on an IP address or localhost", http.StatusMisdirectedRequest)
			return
		}
		h := w.Header()
		h.Set("Content-Security-Policy", contentSecurityPolicy)
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("Referrer-Policy", "no-referrer")
		protected.ServeHTTP(w, r)
	})
}

func localHost(hostport string) bool {
	host, _, err := net.SplitHostPort(hostport)
	if err != nil {
		host = hostport
	}
	host = strings.Trim(host, "[]")
	return strings.EqualFold(host, "localhost") || net.ParseIP(host) != nil
}

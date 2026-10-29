package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGuard(t *testing.T) {
	ok := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNoContent) })
	handler := guard(ok)

	tests := []struct {
		name   string
		method string
		host   string
		origin string
		want   int
	}{
		{"dashboard write", http.MethodPost, "127.0.0.1:7878", "http://127.0.0.1:7878", http.StatusNoContent},
		{"localhost read", http.MethodGet, "localhost:7878", "", http.StatusNoContent},
		{"lan address", http.MethodGet, "192.168.1.20:7878", "", http.StatusNoContent},
		{"other site write", http.MethodPost, "127.0.0.1:7878", "https://evil.example", http.StatusForbidden},
		{"rebound host", http.MethodGet, "evil.example:7878", "", http.StatusMisdirectedRequest},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(tt.method, "/api/jobs", nil)
			r.Host = tt.host
			if tt.origin != "" {
				r.Header.Set("Origin", tt.origin)
				r.Header.Set("Sec-Fetch-Site", "cross-site")
				if tt.origin == "http://"+tt.host {
					r.Header.Set("Sec-Fetch-Site", "same-origin")
				}
			}
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, r)
			if w.Code != tt.want {
				t.Errorf("status = %d, want %d", w.Code, tt.want)
			}
		})
	}
}

package middlewares

import (
	"net"
	"net/http"

	"go.uber.org/zap"
)

func CheckIPAddress(trustedSubnet string, logger zap.SugaredLogger) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if trustedSubnet == "" {
				h.ServeHTTP(w, r)
				return
			}

			ips := r.Header.Get("X-Real-IP")
			if ips == "" {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			ip := net.ParseIP(ips)
			if ip == nil {
				logger.Errorf("не удалось распарсить IP %s", ips)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			_, ipNet, err := net.ParseCIDR(trustedSubnet)
			if err != nil {
				logger.Errorf("не удалось распарсить CIDR %s", trustedSubnet)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if !ipNet.Contains(ip) {
				w.WriteHeader(http.StatusForbidden)
				return
			}
		})
	}
}

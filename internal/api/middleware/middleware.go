package middleware

import (
	"github.com/IakimenkoD/xm-companies-service/internal/service"
	"net"
	"net/http"
	"strings"
)

// Auth

const allowedLocation = "CY"

// IpAddress
func CheckIPAddress(ipChecker service.IpApi) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ip := getUserIP(r)
			location, err := ipChecker.GetUserLocation(ctx, ip)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if location != allowedLocation {
				http.Error(w, "your location unallowed", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func getUserIP(r *http.Request) string {
	ipAddress := r.Header.Get("X-Real-Ip")
	if ipAddress == "" {
		ipAddress = r.Header.Get("X-Forwarded-For")
	}
	if ipAddress == "" {
		ipAddress = r.RemoteAddr
	}
	if len(strings.Split(ipAddress, ":")) > 1 {
		ipAddress, _, _ = net.SplitHostPort(ipAddress)
	}

	return ipAddress
}

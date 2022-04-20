package middleware

import (
	"net/http"
)

// Auth

// IpAddress
func CheckIPAddress(defaultIP string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			//ip, _, _ := net.SplitHostPort(r.RemoteAddr)
			//
			//if !isIPv4(ip) {
			//	ip = defaultIP
			//}
			////TODO

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

//func isIPv4(str string) bool {
//	ip := net.ParseIP(str)
//	return ip != nil && strings.Contains(str, ".")
//}

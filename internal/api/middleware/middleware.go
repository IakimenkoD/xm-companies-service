package middleware

import (
	"github.com/IakimenkoD/xm-companies-service/internal/model"
	"github.com/IakimenkoD/xm-companies-service/internal/service"
	"github.com/dgrijalva/jwt-go"
	"net"
	"net/http"
	"strings"
)

const allowedLocation = "CY"

func CheckAuth(jwtKey []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			c, err := r.Cookie("token")
			if err != nil {
				if err == http.ErrNoCookie {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			//TODO from header also

			tokenStr := c.Value

			claims := &model.Claims{}

			tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
				return jwtKey, nil
			})
			if err != nil {
				if err == jwt.ErrSignatureInvalid {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if !tkn.Valid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func CheckIPAddress(ipChecker service.IpChecker) func(http.Handler) http.Handler {
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
				http.Error(w, "your location is not allowed", http.StatusForbidden)
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

package ratelimit

import (
	"kakigori-api/internal/common"
	"kakigori-api/internal/handler/response"
	"net/http"
	"time"
)

type IpRateLimiter struct {
	ipBlocker *IPBlocker
}

func NewIpRateLimiter(maxRequests int, window, blockDuration time.Duration) *IpRateLimiter {
	return &IpRateLimiter{
		ipBlocker: NewIPBlocker(maxRequests, window, blockDuration),
	}
}

func (l *IpRateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/health" {
			next.ServeHTTP(w, r)
			return
		}

		clientIP := common.GetClientIP(r)

		if blocked, message := l.ipBlocker.IsBlocked(clientIP); blocked {
			response.WriteError(w, http.StatusForbidden, message)
			return
		}

		ok, err := l.ipBlocker.IncrementAndCheck(clientIP)
		if !ok {
			if err != nil {
				response.WriteError(w, http.StatusForbidden, err.Error())
			} else {
				response.WriteError(w, http.StatusForbidden, "Rate limit exceeded")
			}
			return
		}

		next.ServeHTTP(w, r)
	})
}

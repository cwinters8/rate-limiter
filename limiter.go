package ratelimiter

import (
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

type Limiter struct {
	Interval     int // Interval in milliseconds
	CallsAllowed int
	limiter      *rate.Limiter
}

// Limit is a middleware function that will limit the number of calls per interval as defined by the limiter
func (l *Limiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowed := l.limiter.Allow()
		if !allowed {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// NewLimiter creates a new rate limiter with an interval in milliseconds and the number of calls allowed per interval
func NewLimiter(interval int, callsAllowed int) *Limiter {
	limit := rate.Every(time.Duration(interval) * time.Millisecond)
	rateLimiter := rate.NewLimiter(limit, callsAllowed)

	return &Limiter{
		Interval:     interval,
		CallsAllowed: callsAllowed,
		limiter:      rateLimiter,
	}
}

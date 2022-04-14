package ratelimiter

import (
	"fmt"
	"net/http"

	"golang.org/x/time/rate"
)

type Limiter struct {
	Interval     int // Interval in milliseconds
	CallsAllowed int
	limiter      *rate.Limiter
}

func (l *Limiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowed := l.limiter.Allow()
		fmt.Println("allowed:", allowed)
		if !allowed {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func NewLimiter(interval int, callsAllowed int) *Limiter {
	limit := rate.Limit(interval / 1000)
	fmt.Println("limit in seconds:", limit)
	rateLimiter := rate.NewLimiter(limit, callsAllowed)
	return &Limiter{
		Interval:     interval,
		CallsAllowed: callsAllowed,
		limiter:      rateLimiter,
	}
}

package main

import (
	"net/http"

	"golang.org/x/time/rate"
)

type Limiter struct {
	Interval     int // Interval in milliseconds
	CallsAllowed int
	Handler      http.Handler
}

func (l *Limiter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	limiter := rate.NewLimiter(rate.Limit(l.CallsAllowed), l.Interval/1000)
	if !limiter.Allow() {
		http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
		return
	}
	l.Handler.ServeHTTP(w, r)
}

func NewLimiter(interval int, callsAllowed int, handler http.Handler) *Limiter {
	return &Limiter{
		Interval:     interval,
		CallsAllowed: callsAllowed,
		Handler:      handler,
	}
}
